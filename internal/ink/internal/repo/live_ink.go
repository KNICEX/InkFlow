package repo

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
)

var (
	ErrLiveInkNotFound = dao.ErrLiveInkNotFound
)

type LiveInkRepo interface {
	Save(ctx context.Context, ink domain.Ink) (int64, error)
	UpdateStatus(ctx context.Context, ink domain.Ink) error
	FindById(ctx context.Context, id int64) (domain.Ink, error)
	FindByIdAndStatus(ctx context.Context, id int64, status domain.Status) (domain.Ink, error)
	ListByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error)
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
	go func() {
		// 更新了文章，删除首页缓存
		er := repo.cache.DelFirstPage(ctx, ink.Author.Id)
		if er != nil {
			repo.l.WithCtx(ctx).Error("del first page cache error", logx.Error(er),
				logx.Int64("inkId", ink.Id),
				logx.Int64("authorId", ink.Author.Id))
		}
	}()
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
		}()
	case domain.InkStatusPublished:
		go func() {
			er := repo.cache.Set(ctx, ink)
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

func (repo *CachedLiveInkRepo) FindById(ctx context.Context, id int64) (domain.Ink, error) {
	ink, err := repo.cache.Get(ctx, id)
	if err == nil {
		return ink, nil
	}
	entity, err := repo.dao.FindById(ctx, id)
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

func (repo *CachedLiveInkRepo) FindByIdAndStatus(ctx context.Context, id int64, status domain.Status) (domain.Ink, error) {
	if status == domain.InkStatusPublished {
		ink, err := repo.cache.Get(ctx, id)
		if err == nil {
			return ink, nil
		}
	}

	entity, err := repo.dao.FindByIdAndStatus(ctx, id, status.ToInt())
	if err != nil {
		return domain.Ink{}, err
	}
	ink := repo.entityToDomain(entity)
	if status == domain.InkStatusPublished {
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

func (repo *CachedLiveInkRepo) ListByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error) {
	// TODO 这里要不要添加缓存
	inks, err := repo.dao.FindByAuthorId(ctx, authorId, offset, limit)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *CachedLiveInkRepo) FindAll(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	inks, err := repo.dao.FindAll(ctx, maxId, limit)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *CachedLiveInkRepo) FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error) {
	cachedInks, err := repo.cache.GetByIds(ctx, ids)
	if err != nil && !errors.Is(err, cache.ErrKeyNotFound) {
		repo.l.WithCtx(ctx).Error("get ink by ids cache error", logx.Error(err), logx.Any("inkIds", ids))
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
	inks, err := repo.dao.FindByIds(ctx, ids)
	if err != nil {
		repo.l.WithCtx(ctx).Error("find ink by ids from db error", logx.Error(err), logx.Any("ids", ids))
		return nil, err
	}
	for _, ink := range inks {
		cachedInks[ink.Id] = repo.entityToDomain(ink)
	}
	return cachedInks, nil
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
