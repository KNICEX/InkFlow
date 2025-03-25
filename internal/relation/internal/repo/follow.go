package repo

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

var (
	ErrAlreadyFollowed = errors.New("already followed")
)

type FollowRepo interface {
	AddFollowRelation(ctx context.Context, c domain.FollowRelation) error
	RemoveFollowRelation(ctx context.Context, c domain.FollowRelation) error
	FindFollowingList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	FindFlowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	GetFollowStats(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error)
	GetFollowStatBatch(ctx context.Context, uids []int64, viewUid int64) (map[int64]domain.FollowStatistic, error)
}

type CachedFollowRepo struct {
	dao   dao.FollowRelationDAO
	cache cache.FollowCache
	l     logx.Logger
}

func NewCachedFollowRepo(dao dao.FollowRelationDAO, cache cache.FollowCache, l logx.Logger) FollowRepo {
	return &CachedFollowRepo{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *CachedFollowRepo) AddFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := repo.dao.CreateFollowRelation(ctx, repo.toEntity(c))
	if err != nil {
		if errors.Is(err, dao.ErrFollowExist) {
			return ErrAlreadyFollowed
		}
		return err
	}
	err = repo.cache.Follow(ctx, c.FollowerId, c.FolloweeId)
	if err != nil {
		repo.l.WithCtx(ctx).Error("add follow cache error", logx.Error(err), logx.Int64("Uid", c.FollowerId))
	}
	return err
}

func (repo *CachedFollowRepo) RemoveFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := repo.dao.CancelFollow(ctx, repo.toEntity(c))
	if err != nil {
		return err
	}
	err = repo.cache.CancelFollow(ctx, c.FollowerId, c.FolloweeId)
	if err != nil {
		repo.l.WithCtx(ctx).Error("cancel follow cache error", logx.Error(err), logx.Int64("Uid", c.FollowerId))
	}
	return err
}

func (repo *CachedFollowRepo) FindFollowingList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	res, err := repo.dao.FollowList(ctx, uid, maxId, limit)
	if err != nil {
		return nil, err
	}

	self := viewUid == uid
	var followedMap map[int64]bool
	if !self {
		// 不是查看自己的关注列表，查询是否关注
		followedMap, err = repo.dao.FollowedBatch(ctx, viewUid, lo.Map(res, func(item dao.UserFollow, index int) int64 {
			return item.FolloweeId
		}))
		if err != nil {
			repo.l.Error("find followed batch error", logx.Error(err), logx.Int64("uid", uid))
		}
	} else {
		followedMap = make(map[int64]bool)
	}

	return lo.Map(res, func(item dao.UserFollow, index int) domain.FollowInfo {
		var followed bool
		if self {
			followed = true
		} else {
			followed = followedMap[item.FolloweeId]
		}
		return domain.FollowInfo{
			Uid:      item.FolloweeId,
			Followed: followed,
		}
	}), nil
}

func (repo *CachedFollowRepo) FindFlowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	res, err := repo.dao.FollowerList(ctx, uid, maxId, limit)
	if err != nil {
		return nil, err
	}

	// 批量查询是否关注
	var followedMap map[int64]bool
	followedMap, err = repo.dao.FollowedBatch(ctx, viewUid, lo.Map(res, func(item dao.UserFollow, index int) int64 {
		return item.FollowerId
	}))
	if err != nil {
		followedMap = make(map[int64]bool)
		repo.l.Error("find followed batch error", logx.Error(err), logx.Int64("uid", uid))
	} else {
		// 防止nil map
		followedMap = make(map[int64]bool)
	}

	return lo.Map(res, func(item dao.UserFollow, index int) domain.FollowInfo {
		return domain.FollowInfo{
			Uid:      item.FollowerId,
			Followed: followedMap[item.FollowerId],
		}
	}), nil
}

func (repo *CachedFollowRepo) GetFollowStats(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error) {
	res, err := repo.cache.GetStatistic(ctx, uid)
	if err == nil {
		if uid == viewUid {
			// 如果是查看自己的关注数，直接返回
			return res, nil
		}
		followed, er := repo.dao.Followed(ctx, viewUid, uid)
		if er != nil {
			repo.l.Error("get follow statistic cache error", logx.Error(er), logx.Int64("Uid", uid))
			return res, er
		}
		res.Followed = followed
		return res, nil
	}

	if !errors.Is(err, cache.ErrKeyNotFound) {
		repo.l.Error("get follow statistic cache error", logx.Error(err), logx.Int64("Uid", uid))
	}

	var followStats dao.FollowStats
	var followed bool
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		followStats, err = repo.dao.FindFollowStats(ctx, uid)
		return er
	})
	if viewUid != uid {
		eg.Go(func() error {
			var er error
			followed, er = repo.dao.Followed(ctx, viewUid, uid)
			return er
		})
	}
	if err = eg.Wait(); err != nil {
		return domain.FollowStatistic{}, err
	}

	go func() {
		if er := repo.cache.SetStatistic(ctx, domain.FollowStatistic{
			Uid:       uid,
			Followers: followStats.Followers,
			Following: followStats.Following,
		}); er != nil {
			repo.l.Error("set follow statistic cache error", logx.Error(err), logx.Int64("uid", uid))
		}
	}()
	return domain.FollowStatistic{
		Followers: followStats.Followers,
		Following: followStats.Following,
		Followed:  followed,
	}, nil
}

func (repo *CachedFollowRepo) GetFollowStatBatch(ctx context.Context, uids []int64, viewUid int64) (map[int64]domain.FollowStatistic, error) {
	var followedMap map[int64]bool

	eg := errgroup.Group{}
	// 查询是否关注
	eg.Go(func() error {
		var er error
		followedMap, er = repo.dao.FollowedBatch(ctx, viewUid, uids)
		return er
	})

	cachedStatsMap, err := repo.cache.GetStatisticBatch(ctx, uids)
	if err != nil {
		repo.l.Error("get follow statistic batch cache error", logx.Error(err), logx.Any("UserIds", uids))
	}
	if len(cachedStatsMap) == len(uids) {
		// 如果缓存命中，直接返回
		if err = eg.Wait(); err != nil {
			return nil, err
		}
		for _, uid := range uids {
			if _, ok := cachedStatsMap[uid]; !ok {
				cachedStatsMap[uid] = domain.FollowStatistic{
					Followers: cachedStatsMap[uid].Followers,
					Following: cachedStatsMap[uid].Following,
					Followed:  followedMap[uid],
				}
			}
		}
		return cachedStatsMap, nil
	}

	if len(cachedStatsMap) > 0 {
		// 过滤掉已经命中的缓存
		uids = lo.Reject(uids, func(id int64, idx int) bool {
			_, ok := cachedStatsMap[id]
			return ok
		})
	} else {
		cachedStatsMap = make(map[int64]domain.FollowStatistic, len(uids))
	}

	var followStatsMap map[int64]dao.FollowStats
	eg.Go(func() error {
		var er error
		followStatsMap, er = repo.dao.FindFollowStatsBatch(ctx, uids)
		return er
	})

	if err = eg.Wait(); err != nil {
		return nil, err
	}

	for uid, stats := range cachedStatsMap {
		stats.Followed = followedMap[uid]
		cachedStatsMap[uid] = stats
	}

	for uid, stats := range followStatsMap {
		cachedStatsMap[uid] = domain.FollowStatistic{
			Followers: stats.Followers,
			Following: stats.Following,
			Followed:  followedMap[uid],
		}
	}

	go func() {
		// 缓存未命中的数据
		stats := make([]domain.FollowStatistic, 0, len(followStatsMap))
		for _, stat := range followStatsMap {
			stats = append(stats, repo.statsToDomain(stat))
		}
		if er := repo.cache.SetStatisticBatch(ctx, stats); er != nil {
			repo.l.Error("set follow statistic batch cache error", logx.Error(er))
		}
	}()

	return cachedStatsMap, nil
}

func (repo *CachedFollowRepo) statsToDomain(stats dao.FollowStats) domain.FollowStatistic {
	return domain.FollowStatistic{
		Uid:       stats.UserId,
		Followers: stats.Followers,
		Following: stats.Following,
	}
}

func (repo *CachedFollowRepo) toDomain(follow dao.UserFollow) domain.FollowRelation {
	return domain.FollowRelation{
		FollowerId: follow.FollowerId,
		FolloweeId: follow.FolloweeId,
		CreatedAt:  follow.CreatedAt,
	}
}
func (repo *CachedFollowRepo) toEntity(follow domain.FollowRelation) dao.UserFollow {
	return dao.UserFollow{
		FollowerId: follow.FollowerId,
		FolloweeId: follow.FolloweeId,
		CreatedAt:  follow.CreatedAt,
	}
}
