package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
)

type InteractiveRepo interface {
	CreateInteractive(ctx context.Context, biz string, bizId int64) error
	IncrView(ctx context.Context, biz string, bizId, uid int64) error
	IncrViewBatch(ctx context.Context, biz string, bizIds, uid []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	IncrReply(ctx context.Context, biz string, bizId int64) error
	DecrReply(ctx context.Context, biz string, bizId int64) error
	IncrFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error
	DecrFavorite(ctx context.Context, biz string, bizId, uid int64) error

	ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error)
	ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error)
	ListFavoriteByFid(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]domain.FavoriteRecord, error)

	Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	LikedBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]bool, error)
	Favorited(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	FavoritedBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]bool, error)
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	GetBatch(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
	CountUserLikes(ctx context.Context, uid int64) (int64, error)
	CountUserViews(ctx context.Context, uid int64) (int64, error)
}

type CachedInteractiveRepo struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logx.Logger
}

func NewCachedInteractiveRepo(cache cache.InteractiveCache, dao dao.InteractiveDAO, l logx.Logger) InteractiveRepo {
	return &CachedInteractiveRepo{
		cache: cache,
		dao:   dao,
		l:     l,
	}
}

func (repo *CachedInteractiveRepo) CreateInteractive(ctx context.Context, biz string, bizId int64) error {
	return repo.dao.InsertInteractive(ctx, biz, bizId)
}

func (repo *CachedInteractiveRepo) IncrView(ctx context.Context, biz string, bizId, uid int64) error {
	if err := repo.dao.InsertView(ctx, biz, bizId, uid); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.IncrViewCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("incr read cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) IncrViewBatch(ctx context.Context, biz string, bizIds, uids []int64) error {
	if err := repo.dao.InsertViewBatch(ctx, biz, bizIds, uids); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.IncrViewCntBatch(context.WithoutCancel(ctx), biz, bizIds); err != nil {
			repo.l.WithCtx(ctx).Error("incr read cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Any("bizIds", bizIds),
				logx.Any("uids", uids))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	if err := repo.dao.InsertLike(ctx, biz, bizId, uid); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.IncrLikeCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("incr like cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	if err := repo.dao.DeleteLike(ctx, biz, bizId, uid); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.DecrLikeCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("decr like cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) IncrReply(ctx context.Context, biz string, bizId int64) error {
	if err := repo.dao.IncrReply(ctx, biz, bizId); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.IncrReplyCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("incr reply cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) DecrReply(ctx context.Context, biz string, bizId int64) error {
	if err := repo.dao.DecrReply(ctx, biz, bizId); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.DecrReplyCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("decr reply cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) IncrFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error {
	if err := repo.dao.InsertFavorite(ctx, biz, bizId, uid, fid); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.IncrFavoriteCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("incr favorite cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) DecrFavorite(ctx context.Context, biz string, bizId, uid int64) error {
	if err := repo.dao.DeleteFavorite(ctx, biz, bizId, uid); err != nil {
		return err
	}
	go func() {
		if err := repo.cache.DecrFavoriteCnt(context.WithoutCancel(ctx), biz, bizId); err != nil {
			repo.l.WithCtx(ctx).Error("decr favorite cnt cache error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}()
	return nil
}

func (repo *CachedInteractiveRepo) ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error) {
	records, err := repo.dao.ListViewRecord(ctx, biz, uid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(records, func(item dao.UserView, index int) domain.ViewRecord {
		return repo.readToDomain(item)
	}), nil
}

func (repo *CachedInteractiveRepo) ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error) {
	records, err := repo.dao.ListLikeRecord(ctx, biz, uid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(records, func(item dao.UserLike, index int) domain.LikeRecord {
		return repo.likeToDomain(item)
	}), nil
}

func (repo *CachedInteractiveRepo) ListFavoriteByFid(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]domain.FavoriteRecord, error) {
	records, err := repo.dao.FindByFavorite(ctx, biz, uid, fid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(records, func(item dao.UserFavorite, index int) domain.FavoriteRecord {
		return repo.favoriteToDomain(item)

	}), nil
}

func (repo *CachedInteractiveRepo) Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	liked, err := repo.dao.FindLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return false, err
	}
	return liked.Id != 0, nil
}

func (repo *CachedInteractiveRepo) LikedBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]bool, error) {
	liked, err := repo.dao.FindLikeBatch(ctx, biz, bizIds, uid)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]bool)
	for _, item := range liked {
		res[item.BizId] = true
	}
	return res, nil
}

func (repo *CachedInteractiveRepo) Favorited(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	favorited, err := repo.dao.FindFavoriteInfo(ctx, biz, bizId, uid)
	if err != nil {
		return false, err
	}
	return favorited.Id != 0, nil
}

func (repo *CachedInteractiveRepo) FavoritedBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]bool, error) {
	favorited, err := repo.dao.FindFavoriteBatch(ctx, biz, bizIds, uid)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]bool)
	for _, item := range favorited {
		res[item.BizId] = true
	}
	return res, nil
}

func (repo *CachedInteractiveRepo) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := repo.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}

	entity, err := repo.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = repo.interToDomain(entity)

	go func() {
		if er := repo.cache.Set(context.WithoutCancel(ctx), biz, bizId, intr); er != nil {
			repo.l.WithCtx(ctx).Error("set interactive cache error", logx.Error(er),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId))
		}
	}()
	return intr, nil
}

func (repo *CachedInteractiveRepo) GetBatch(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	intrMap, err := repo.cache.GetBatch(ctx, biz, bizIds)
	if err != nil {
		repo.l.WithCtx(ctx).Error("get interactive batch cache error", logx.Error(err),
			logx.String("biz", biz),
			logx.Any("bizIds", bizIds))
	}

	if len(intrMap) == len(bizIds) {
		return intrMap, nil
	}

	if len(intrMap) > 0 {
		// 过滤掉命中缓存的 id
		bizIds = lo.Reject(bizIds, func(id int64, idx int) bool {
			_, ok := intrMap[id]
			return ok
		})
	} else {
		intrMap = make(map[int64]domain.Interactive, len(bizIds))
	}

	entities, err := repo.dao.GetByIds(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}
	intrMap = make(map[int64]domain.Interactive)
	for _, entity := range entities {
		intrMap[entity.BizId] = repo.interToDomain(entity)
	}

	// 提前转切片，避免并发读写map
	intrs := lo.MapToSlice(intrMap, func(key int64, value domain.Interactive) domain.Interactive {
		return value
	})
	go func() {
		if er := repo.cache.SetBatch(context.WithoutCancel(ctx), intrs); er != nil {
			repo.l.WithCtx(ctx).Error("set interactive cache error", logx.Error(er),
				logx.String("biz", biz),
				logx.Any("bizIds", bizIds))
		}
	}()
	return intrMap, nil
}

func (repo *CachedInteractiveRepo) CountUserLikes(ctx context.Context, uid int64) (int64, error) {
	return repo.dao.CountUserLikes(ctx, uid)
}

func (repo *CachedInteractiveRepo) CountUserViews(ctx context.Context, uid int64) (int64, error) {
	return repo.dao.CountUserViews(ctx, uid)
}

func (repo *CachedInteractiveRepo) interToDomain(entity dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz:     entity.Biz,
		BizId:   entity.BizId,
		ViewCnt: entity.ViewCnt,
		LikeCnt: entity.LikeCnt,
	}
}

func (repo *CachedInteractiveRepo) readToDomain(entity dao.UserView) domain.ViewRecord {
	return domain.ViewRecord{
		Biz:    entity.Biz,
		BizId:  entity.BizId,
		UserId: entity.UserId,
	}
}
func (repo *CachedInteractiveRepo) likeToDomain(entity dao.UserLike) domain.LikeRecord {
	return domain.LikeRecord{
		Biz:    entity.Biz,
		BizId:  entity.BizId,
		UserId: entity.UserId,
	}
}

func (repo *CachedInteractiveRepo) favoriteToDomain(entity dao.UserFavorite) domain.FavoriteRecord {
	return domain.FavoriteRecord{
		Biz:    entity.Biz,
		BizId:  entity.BizId,
		Fid:    entity.FavoriteId,
		UserId: entity.UserId,
	}
}
