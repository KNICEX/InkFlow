package cache

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var ErrKeyNotFound = redis.Nil

type FollowCache interface {
	GetStatisticInfo(ctx context.Context, uid int64) (domain.FollowStatistic, error)
	SetStatisticInfo(ctx context.Context, uid int64, statistic domain.FollowStatistic) error
	Follow(ctx context.Context, uid, followeeId int64) error
	CancelFollow(ctx context.Context, uid, followeeId int64) error
}

type RedisFollowCache struct {
	cmd redis.Cmdable
}

const (
	fieldFollowerCount = "follower_count"
	fieldFolloweeCount = "followee_count"
)

func (c *RedisFollowCache) statisticKey(uid int64) string {
	return fmt.Sprintf("relation:follow:statistic:%d", uid)
}
func (c *RedisFollowCache) GetStatisticInfo(ctx context.Context, uid int64) (domain.FollowStatistic, error) {
	data, err := c.cmd.HGetAll(ctx, c.statisticKey(uid)).Result()
	res := domain.FollowStatistic{}
	if err != nil {
		return res, err
	}
	if len(data) == 0 {
		return res, ErrKeyNotFound
	}
	res.Followers, _ = strconv.ParseInt(data[fieldFollowerCount], 10, 64)
	res.Followings, _ = strconv.ParseInt(data[fieldFolloweeCount], 10, 64)
	return res, nil
}

func (c *RedisFollowCache) SetStatisticInfo(ctx context.Context, uid int64, statistic domain.FollowStatistic) error {
	return c.cmd.HMSet(ctx, c.statisticKey(uid), fieldFollowerCount, statistic.Followers, fieldFolloweeCount, statistic.Followings).Err()
}

func (c *RedisFollowCache) updateStatisticInfo(ctx context.Context, followerId, followeeId, delta int64) error {
	tx := c.cmd.Pipeline()
	// 关注者的关注数 + delta
	tx.HIncrBy(ctx, c.statisticKey(followerId), fieldFolloweeCount, delta)
	// 被关注者的粉丝数 + delta
	tx.HIncrBy(ctx, c.statisticKey(followeeId), fieldFollowerCount, delta)
	_, err := tx.Exec(ctx)
	return err
}

func (c *RedisFollowCache) Follow(ctx context.Context, uid, followeeId int64) error {
	return c.updateStatisticInfo(ctx, uid, followeeId, 1)
}

func (c *RedisFollowCache) CancelFollow(ctx context.Context, uid, followeeId int64) error {
	return c.updateStatisticInfo(ctx, uid, followeeId, -1)
}
