package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/stringx"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"strings"

	"github.com/KNICEX/InkFlow/internal/comment/internal/domain"
)

// CommentRepo defines the data access operations for comments
type CommentRepo interface {
	CreateComment(ctx context.Context, comment domain.Comment) (int64, error)
	DelComment(ctx context.Context, id int64) error
	LikeComment(ctx context.Context, uid, cid int64) error
	CancelLike(ctx context.Context, uid, cid int64) error

	FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]domain.Comment, error)
	FindByRootId(ctx context.Context, rootId int64, maxId int64, limit int) ([]domain.Comment, error)
	FindByParentId(ctx context.Context, parentId int64, maxId int64, limit int) ([]domain.Comment, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Comment, error)
	FindById(ctx context.Context, id int64) (domain.Comment, error)
	FindAuthorReplyIn(ctx context.Context, ids []int64) (map[int64][]domain.Comment, error)

	FindStats(ctx context.Context, ids []int64, uid int64) (map[int64]domain.CommentStats, error)
	BizReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error)
}

type CachedCommentRepo struct {
	dao   dao.CommentDAO
	cache cache.CommentCache
	l     logx.Logger
}

func NewCachedCommentRepo(dao dao.CommentDAO, cache cache.CommentCache, l logx.Logger) CommentRepo {
	return &CachedCommentRepo{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *CachedCommentRepo) CreateComment(ctx context.Context, comment domain.Comment) (int64, error) {
	if comment.Root == nil {
		// 一级评论, 增加评论数
		go func() {
			er := repo.cache.IncrBizReply(ctx, comment.Biz, comment.BizId)
			if er != nil {
				repo.l.WithCtx(ctx).Error("comment cache incr biz reply error",
					logx.String("biz", comment.Biz),
					logx.Int64("bizId", comment.BizId),
					logx.Error(er))
			}
		}()
	}
	return repo.dao.Insert(ctx, repo.toEntity(comment))
}

func (repo *CachedCommentRepo) DelComment(ctx context.Context, id int64) error {
	// TODO 这里需要减少对应biz的评论数缓存
	return repo.dao.Delete(ctx, id)
}

func (repo *CachedCommentRepo) LikeComment(ctx context.Context, uid, cid int64) error {
	return repo.dao.Like(ctx, uid, cid)
}

func (repo *CachedCommentRepo) CancelLike(ctx context.Context, uid, cid int64) error {
	return repo.dao.CancelLike(ctx, uid, cid)
}

func (repo *CachedCommentRepo) FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := repo.dao.FindByBiz(ctx, biz, bizId, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(comments, func(item dao.Comment, index int) domain.Comment {
		return repo.toDomain(item)
	}), nil
}

func (repo *CachedCommentRepo) FindByRootId(ctx context.Context, rootId int64, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := repo.dao.FindRepliesByRid(ctx, rootId, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(comments, func(item dao.Comment, index int) domain.Comment {
		return repo.toDomain(item)
	}), nil
}

func (repo *CachedCommentRepo) FindByParentId(ctx context.Context, parentId int64, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := repo.dao.FindRepliesByPid(ctx, parentId, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(comments, func(item dao.Comment, index int) domain.Comment {
		return repo.toDomain(item)
	}), nil
}
func (repo *CachedCommentRepo) FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Comment, error) {
	comments, err := repo.dao.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Comment)
	for _, item := range comments {
		res[item.Id] = repo.toDomain(item)
	}
	return res, nil
}
func (repo *CachedCommentRepo) FindById(ctx context.Context, id int64) (domain.Comment, error) {
	comment, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.Comment{}, err
	}
	return repo.toDomain(comment), nil
}

func (repo *CachedCommentRepo) FindAuthorReplyIn(ctx context.Context, ids []int64) (map[int64][]domain.Comment, error) {
	comments, err := repo.dao.FindAuthorReplyIn(ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64][]domain.Comment)
	for key, replies := range comments {
		res[key] = lo.Map(replies, func(item dao.Comment, index int) domain.Comment {
			return repo.toDomain(item)
		})
	}
	return res, nil
}
func (repo *CachedCommentRepo) BizReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error) {
	cachedCounts, err := repo.cache.BizReplyCount(ctx, biz, bizIds)
	if err != nil {
		repo.l.WithCtx(ctx).Error("comment cache get biz reply count error",
			logx.String("biz", biz),
			logx.Any("bizIds", bizIds),
			logx.Error(err))
	}
	if len(cachedCounts) == len(bizIds) {
		// 全部命中缓存
		return cachedCounts, nil
	}

	if len(cachedCounts) > 0 {
		// 过滤掉已经缓存的bizIds
		bizIds = lo.Reject(bizIds, func(item int64, index int) bool {
			_, ok := cachedCounts[item]
			return ok
		})
	}

	counts, err := repo.dao.ReplyCount(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	go func() {
		er := repo.cache.SetBizReplyCount(ctx, biz, counts)
		if er != nil {
			repo.l.WithCtx(ctx).Error("comment cache set biz reply count error",
				logx.String("biz", biz),
				logx.Any("bizIds", bizIds),
				logx.Error(er))
		}
	}()

	for id, cnt := range cachedCounts {
		counts[id] = cnt
	}
	return counts, nil
}

func (repo *CachedCommentRepo) FindStats(ctx context.Context, ids []int64, uid int64) (map[int64]domain.CommentStats, error) {
	var stats map[int64]dao.CommentStats
	var likedMap map[int64]bool
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		stats, err = repo.dao.FindStats(ctx, ids)
		return err
	})
	eg.Go(func() error {
		var err error
		likedMap, err = repo.dao.Liked(ctx, uid, ids)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	res := make(map[int64]domain.CommentStats)
	for _, item := range stats {
		res[item.CommentId] = domain.CommentStats{
			LikeCnt:  item.LikeCount,
			ReplyCnt: item.ReplyCount,
			Liked:    likedMap[item.CommentId],
		}
	}
	return res, nil
}

func (repo *CachedCommentRepo) toDomain(en dao.Comment) domain.Comment {
	var root, parent *domain.Comment
	if en.RootId > 0 {
		root = &domain.Comment{
			Id: en.RootId,
		}
	}
	if en.ParentId > 0 {
		parent = &domain.Comment{
			Id: en.ParentId,
		}
	}
	return domain.Comment{
		Id:    en.Id,
		Biz:   en.Biz,
		BizId: en.BizId,
		Commentator: domain.Commentator{
			Id:       en.CommentatorId,
			IsAuthor: en.IsAuthor,
		},
		Root:   root,
		Parent: parent,
		Payload: domain.Payload{
			Content: en.Content,
			Images:  stringx.Split(en.Images, ","),
		},
		CreatedAt: en.CreatedAt,
	}
}

func (repo *CachedCommentRepo) toEntity(comment domain.Comment) dao.Comment {
	var rootId, parentId int64
	if comment.Root != nil {
		rootId = comment.Root.Id
	}
	if comment.Parent != nil {
		parentId = comment.Parent.Id
	}
	return dao.Comment{
		Id:            comment.Id,
		Biz:           comment.Biz,
		BizId:         comment.BizId,
		RootId:        rootId,
		ParentId:      parentId,
		Content:       comment.Payload.Content,
		Images:        strings.Join(comment.Payload.Images, ","),
		CommentatorId: comment.Commentator.Id,
		IsAuthor:      comment.Commentator.IsAuthor,
		CreatedAt:     comment.CreatedAt,
	}
}
