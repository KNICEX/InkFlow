package repo

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
	"slices"
)

var (
	ErrLiveInkNotFound = dao.ErrLiveInkNotFound
)

type LiveInkRepo interface {
	Save(ctx context.Context, ink domain.Ink) (int64, error)
	UpdateStatus(ctx context.Context, ink domain.Ink) error
	UpdateAiTags(ctx context.Context, id int64, tags domain.Tags) error
	Delete(ctx context.Context, id int64, authorId int64, status ...domain.Status) error
	FindById(ctx context.Context, id int64, status ...domain.Status) (domain.Ink, error)
	FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...domain.Status) ([]domain.Ink, error)
	FindAll(ctx context.Context, maxId int64, limit int, status ...domain.Status) ([]domain.Ink, error)
	FindByIds(ctx context.Context, ids []int64, status ...domain.Status) (map[int64]domain.Ink, error)
}

// CachedLiveInkRepo
// 考虑在这里做一个命中率统计
type CachedLiveInkRepo struct {
	dao   dao.LiveDAO
	cache cache.InkCache
	l     logx.Logger
}

func NewCachedLiveInkRepo(dao dao.LiveDAO, cache cache.InkCache, l logx.Logger) LiveInkRepo {
	return &CachedLiveInkRepo{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *CachedLiveInkRepo) Save(ctx context.Context, ink domain.Ink) (int64, error) {
	id, err := repo.dao.Upsert(ctx, repo.domainToEntity(ink))
	if err != nil {
		return 0, err
	}
	go func() {
		// 更新了文章，删除首页缓存
		er := repo.cache.DelFirstPage(ctx, ink.Author.Id)
		if er != nil {
			repo.l.WithCtx(ctx).Error("del first page cache error", logx.Error(er),
				logx.Int64("inkId", ink.Id),
				logx.Int64("authorId", ink.Author.Id))
		}
	}()
	if ink.Status == domain.InkStatusPublished {
		go func() {
			// 新文章查看概率高，缓存时间短
			er := repo.cache.Set(ctx, ink)
			if err != nil {
				repo.l.WithCtx(ctx).Error("set ink cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}
		}()
	}
	return id, nil
}

func (repo *CachedLiveInkRepo) UpdateStatus(ctx context.Context, ink domain.Ink) error {
	err := repo.dao.UpdateStatus(ctx, ink.Id, ink.Author.Id, ink.Status.ToInt())
	if err != nil {
		return err
	}
	switch ink.Status {
	case domain.InkStatusPrivate, domain.InkStatusUnPublished:
		// 如果是隐藏了文章或者退回到草稿，删除缓存
		go func() {
			er := repo.cache.Del(ctx, ink.Id)
			if err != nil {
				repo.l.WithCtx(ctx).Error("del ink cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}
			// 删除首页缓存
			er = repo.cache.DelFirstPage(ctx, ink.Author.Id)
			if er != nil {
				repo.l.WithCtx(ctx).Error("del first page cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}
		}()
	case domain.InkStatusPublished:
		go func() {
			// 文章发布成功，设置预缓存，并且删除首页缓存
			// 删除首页缓存
			if er := repo.cache.DelFirstPage(ctx, ink.Author.Id); er != nil {
				repo.l.WithCtx(ctx).Error("del first page cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}
			// 设置预缓存
			entity, er := repo.dao.FindById(ctx, ink.Id)
			if er != nil {
				repo.l.WithCtx(ctx).Error("published pre cache find ink by id error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
				return
			}
			er = repo.cache.Set(ctx, repo.entityToDomain(entity))
			if err != nil {
				repo.l.WithCtx(ctx).Error("set ink cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}

		}()
	default:

	}

	return nil
}

func (repo *CachedLiveInkRepo) UpdateAiTags(ctx context.Context, id int64, tags domain.Tags) error {
	err := repo.dao.UpdateAiTags(ctx, id, tags.ToString())
	if err != nil {
		return err
	}
	return nil
}

func (repo *CachedLiveInkRepo) Delete(ctx context.Context, id int64, authorId int64, status ...domain.Status) error {
	err := repo.dao.Delete(ctx, id, authorId, repo.parseStatus(status)...)
	if err != nil {
		return err
	}

	if len(status) == 0 || slices.Contains(status, domain.InkStatusPublished) {
		// 如果是删除了已发布的文章，删除缓存
		go func() {
			er := repo.cache.Del(ctx, id)
			if er != nil {
				repo.l.WithCtx(ctx).Error("del ink cache error", logx.Error(er),
					logx.Int64("inkId", id),
					logx.Int64("authorId", authorId))
			}
			// 删除首页缓存
			er = repo.cache.DelFirstPage(ctx, authorId)
			if er != nil {
				repo.l.WithCtx(ctx).Error("del first page cache error", logx.Error(er),
					logx.Int64("inkId", id),
					logx.Int64("authorId", authorId))
			}
		}()
	}
	return nil
}

func (repo *CachedLiveInkRepo) FindById(ctx context.Context, id int64, status ...domain.Status) (domain.Ink, error) {
	var ink domain.Ink
	var err error
	if len(status) == 0 || slices.Contains(status, domain.InkStatusPublished) {
		ink, err = repo.cache.Get(ctx, id)
		if err == nil {
			return ink, nil
		}
	}

	entity, err := repo.dao.FindById(ctx, id, repo.parseStatus(status)...)
	if err != nil {
		return domain.Ink{}, err
	}
	ink = repo.entityToDomain(entity)
	if ink.Status == domain.InkStatusPublished {
		// 缓存未命中且是已发布的文章，设置缓存
		go func() {
			er := repo.cache.Set(ctx, ink)
			if er != nil {
				repo.l.WithCtx(ctx).Error("set ink cache error", logx.Error(er),
					logx.Int64("inkId", ink.Id),
					logx.Int64("authorId", ink.Author.Id))
			}
		}()
	}
	return ink, nil
}

func (repo *CachedLiveInkRepo) FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...domain.Status) ([]domain.Ink, error) {
	inks, err := repo.dao.FindByAuthorId(ctx, authorId, offset, limit, repo.parseStatus(status)...)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *CachedLiveInkRepo) FindAll(ctx context.Context, maxId int64, limit int, status ...domain.Status) ([]domain.Ink, error) {
	inks, err := repo.dao.FindAll(ctx, maxId, limit, repo.parseStatus(status)...)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *CachedLiveInkRepo) FindByIds(ctx context.Context, ids []int64, status ...domain.Status) (map[int64]domain.Ink, error) {
	var cachedInks map[int64]domain.Ink
	var err error
	if len(ids) == 0 || slices.Contains(status, domain.InkStatusPublished) {
		// 无特定或者已发布状态的文章，才查缓存
		cachedInks, err = repo.cache.GetByIds(ctx, ids)
		if err != nil && !errors.Is(err, cache.ErrKeyNotFound) {
			repo.l.WithCtx(ctx).Error("get ink by ids cache error", logx.Error(err), logx.Any("inkIds", ids))
		}
	}

	if len(cachedInks) == len(ids) {
		return cachedInks, nil
	}

	// 去除缓存中查询到的
	if len(cachedInks) > 0 {
		ids = lo.WithoutBy(ids, func(id int64) bool {
			_, ok := cachedInks[id]
			return ok
		})
	}
	inks, err := repo.dao.FindByIds(ctx, ids, repo.parseStatus(status)...)
	if err != nil {
		repo.l.WithCtx(ctx).Error("find ink by ids from db error", logx.Error(err), logx.Any("ids", ids))
		return nil, err
	}
	for _, ink := range inks {
		cachedInks[ink.Id] = repo.entityToDomain(ink)
	}
	return cachedInks, nil
}

func (repo *CachedLiveInkRepo) parseStatus(status []domain.Status) []int {
	if len(status) == 0 {
		return nil
	}
	return lo.Map(status, func(item domain.Status, index int) int {
		return item.ToInt()
	})
}

func (repo *CachedLiveInkRepo) domainToEntity(ink domain.Ink) dao.LiveInk {
	return dao.LiveInk{
		Id:          ink.Id,
		AuthorId:    ink.Author.Id,
		Title:       ink.Title,
		Cover:       ink.Cover,
		Summary:     ink.Summary,
		CategoryId:  ink.Category.Id,
		ContentType: ink.ContentType.ToInt(),
		Tags:        ink.Tags.ToString(),
		AiTags:      ink.AiTags.ToString(),
		ContentHtml: ink.ContentHtml,
		ContentMeta: ink.ContentMeta,
		Status:      int(ink.Status),
		CreatedAt:   ink.CreatedAt,
		UpdatedAt:   ink.UpdatedAt,
	}
}
func (repo *CachedLiveInkRepo) entityToDomain(ink dao.LiveInk) domain.Ink {
	return domain.Ink{
		Id: ink.Id,
		Author: domain.Author{
			Id: ink.AuthorId,
		},
		Title:   ink.Title,
		Cover:   ink.Cover,
		Summary: ink.Summary,
		Category: domain.Category{
			Id: ink.CategoryId,
		},
		ContentType: domain.ContentTypeFromInt(ink.ContentType),
		Tags:        domain.TagsFromString(ink.Tags),
		AiTags:      domain.TagsFromString(ink.AiTags),
		ContentHtml: ink.ContentHtml,
		ContentMeta: ink.ContentMeta,
		Status:      domain.Status(ink.Status),
		CreatedAt:   ink.CreatedAt,
		UpdatedAt:   ink.UpdatedAt,
	}
}
