package service

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/comment/internal/domain"
	"github.com/KNICEX/InkFlow/internal/comment/internal/event"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	ErrNoPermission = errors.New("no permission")
)

const (
	bizInk = "ink"
)

type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) (int64, error)
	Delete(ctx context.Context, id int64, uid int64) error
	Like(ctx context.Context, uid, cid int64) error

	FindById(ctx context.Context, commentId, uid int64) (domain.Comment, error)
	FindByIds(ctx context.Context, ids []int64, uid int64) (map[int64]domain.Comment, error)

	// LoadLastedList 加载一级评论列表
	LoadLastedList(ctx context.Context, biz string, bizId int64, uid, maxId int64, limit int) ([]domain.Comment, error)
	// LoadMoreRepliesByRid 根据rootId加载所有子评论
	LoadMoreRepliesByRid(ctx context.Context, rid int64, uid, maxId int64, limit int) ([]domain.Comment, error)
	// LoadMoreRepliesByPid 根据parentId加载所有子评论
	LoadMoreRepliesByPid(ctx context.Context, pid int64, uid, maxId int64, limit int) ([]domain.Comment, error)
}

type commentService struct {
	repo     repo.CommentRepo
	l        logx.Logger
	inkSvc   ink.Service
	producer event.CommentEvtProducer
}

func NewCommentService(repo repo.CommentRepo, inkSvc ink.Service, producer event.CommentEvtProducer, l logx.Logger) CommentService {
	return &commentService{
		repo:     repo,
		l:        l,
		inkSvc:   inkSvc,
		producer: producer,
	}
}

func (svc *commentService) LoadLastedList(ctx context.Context, biz string, bizId int64, uid, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := svc.repo.FindByBiz(ctx, biz, bizId, maxId, limit)
	if err != nil {
		return nil, err
	}
	cids := lo.Map(comments, func(item domain.Comment, index int) int64 {
		return item.Id
	})

	authorReplies, err := svc.repo.FindAuthorReplyIn(ctx, cids)
	if err != nil {
		return nil, err
	}
	for _, replies := range authorReplies {
		cids = append(cids, lo.Map(replies, func(item domain.Comment, index int) int64 {
			return item.Id
		})...)
	}
	stats, err := svc.repo.FindStats(ctx, cids, uid)
	if err != nil {
		return nil, err
	}

	for i, comment := range comments {
		authorRe := authorReplies[comment.Id]
		for j, reply := range authorRe {
			// 组装评作者回复的统计数据
			authorRe[j].Stats = stats[reply.Id]
		}
		comment.Children = authorRe
		comment.Stats = stats[comment.Id]
		comments[i] = comment
	}
	return comments, nil
}

func (svc *commentService) Create(ctx context.Context, comment domain.Comment) (int64, error) {
	isAuthor, err := svc.isAuthor(ctx, comment.Biz, comment.BizId, comment.Commentator.Id)
	if err != nil {
		return 0, err
	}
	comment.Commentator.IsAuthor = isAuthor
	id, err := svc.repo.CreateComment(ctx, comment)
	if err != nil {
		return 0, err
	}
	go func() {
		var rootId, parentId int64
		if comment.Root != nil {
			rootId = comment.Root.Id
		}
		if comment.Parent != nil {
			parentId = comment.Parent.Id
		}
		er := svc.producer.ProduceReply(ctx, event.ReplyEvent{
			CommentId:     id,
			RootId:        rootId,
			ParentId:      parentId,
			Biz:           comment.Biz,
			BizId:         comment.BizId,
			CommentatorId: comment.Commentator.Id,
			Payload: event.Payload{
				Content: comment.Payload.Content,
				Images:  comment.Payload.Images,
			},
			CreatedAt: time.Now(),
		})
		if er != nil {
			svc.l.WithCtx(ctx).Error("produce reply event error", logx.Error(er),
				logx.Int64("commentId", id),
				logx.Int64("uid", comment.Commentator.Id))
		}
	}()
	return id, nil
}

func (svc *commentService) isAuthor(ctx context.Context, biz string, bizId int64, uid int64) (bool, error) {
	switch biz {
	case bizInk:
		inkInfo, err := svc.inkSvc.FindById(ctx, bizId)
		if err != nil {
			return false, err
		}
		return inkInfo.Author.Id == uid, nil
	default:
		return false, nil
	}
}

func (svc *commentService) Delete(ctx context.Context, id int64, uid int64) error {
	c, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if c.Commentator.Id != uid {
		return ErrNoPermission
	}

	go func() {
		er := svc.producer.ProduceDelete(ctx, event.DeleteEvent{
			CommentId: id,
			CreatedAt: time.Now(),
		})
		if er != nil {
			svc.l.WithCtx(ctx).Error("produce delete event error", logx.Error(err),
				logx.Int64("commentId", id),
				logx.Int64("uid", uid))
		}
	}()
	return svc.repo.DelComment(ctx, id)
}

func (svc *commentService) Like(ctx context.Context, uid, cid int64) error {
	err := svc.repo.LikeComment(ctx, uid, cid)
	if err != nil {
		return err
	}
	go func() {
		er := svc.producer.ProduceLike(ctx, event.LikeEvent{
			CommentId: cid,
			LikeUid:   uid,
			CreatedAt: time.Now(),
		})
		if er != nil {
			svc.l.WithCtx(ctx).Error("produce like event error", logx.Error(er),
				logx.Int64("commentId", cid),
				logx.Int64("likeUid", uid))
		}
	}()
	return nil
}

func (svc *commentService) FindById(ctx context.Context, commentId, uid int64) (domain.Comment, error) {
	eg := errgroup.Group{}
	var comment domain.Comment
	var stats map[int64]domain.CommentStats
	eg.Go(func() error {
		var err error
		comment, err = svc.repo.FindById(ctx, commentId)
		return err
	})
	eg.Go(func() error {
		var err error
		stats, err = svc.repo.FindStats(ctx, []int64{commentId}, uid)
		return err
	})

	if err := eg.Wait(); err != nil {
		return domain.Comment{}, err
	}
	comment.Stats = stats[comment.Id]
	return comment, nil
}

func (svc *commentService) FindByIds(ctx context.Context, ids []int64, uid int64) (map[int64]domain.Comment, error) {
	comments, err := svc.repo.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	stats, err := svc.repo.FindStats(ctx, ids, uid)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		comment.Stats = stats[comment.Id]
	}
	return comments, nil
}

func (svc *commentService) LoadMoreRepliesByRid(ctx context.Context, rid int64, uid, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := svc.repo.FindByRootId(ctx, rid, maxId, limit)
	if err != nil {
		return nil, err
	}
	cids := lo.Map(comments, func(item domain.Comment, index int) int64 {
		return item.Id
	})

	stats, err := svc.repo.FindStats(ctx, cids, uid)
	if err != nil {
		return nil, err
	}

	for i, comment := range comments {
		comment.Stats = stats[comment.Id]
		comments[i] = comment
	}
	return comments, nil
}

func (svc *commentService) LoadMoreRepliesByPid(ctx context.Context, pid int64, uid, maxId int64, limit int) ([]domain.Comment, error) {
	comments, err := svc.repo.FindByParentId(ctx, pid, maxId, limit)
	if err != nil {
		return nil, err
	}
	cids := lo.Map(comments, func(item domain.Comment, index int) int64 {
		return item.Id
	})

	stats, err := svc.repo.FindStats(ctx, cids, uid)
	if err != nil {
		return nil, err
	}

	for i, comment := range comments {
		comment.Stats = stats[comment.Id]
		comments[i] = comment
	}
	return comments, nil
}
