package cache

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var ErrKeyNotFound = redis.Nil

type FollowCache interface {
	GetStatistic(ctx context.Context, uid int64) (domain.FollowStatistic, error)
	GetStatisticBatch(ctx context.Context, uids []int64) (map[int64]domain.FollowStatistic, error)
	SetStatistic(ctx context.Context, stats domain.FollowStatistic) error
	SetStatisticBatch(ctx context.Context, stats []domain.FollowStatistic) error
	Follow(ctx context.Context, uid, followeeId int64) error
	CancelFollow(ctx context.Context, uid, followeeId int64) error
}

type RedisFollowCache struct {
	cmd redis.Cmdable
	exp time.Duration
}

func NewRedisFollowCache(cmd redis.Cmdable) FollowCache {
	return &RedisFollowCache{
		cmd: cmd,
		exp: time.Minute * 10,
	}
}

const (
	fieldFollowerCount = "follower_count"
	fieldFolloweeCount = "followee_count"
)

func (c *RedisFollowCache) statisticKey(uid int64) string {
	return fmt.Sprintf("relation:follow:statistic:%d", uid)
}
func (c *RedisFollowCache) GetStatistic(ctx context.Context, uid int64) (domain.FollowStatistic, error) {
	data, err := c.cmd.HGetAll(ctx, c.statisticKey(uid)).Result()
	res := domain.FollowStatistic{}
	if err != nil {
		return res, err
	}
	if len(data) == 0 {
		return res, ErrKeyNotFound
	}
	res.Followers, _ = strconv.ParseInt(data[fieldFollowerCount], 10, 64)
	res.Following, _ = strconv.ParseInt(data[fieldFolloweeCount], 10, 64)
	return res, nil
}

func (c *RedisFollowCache) GetStatisticBatch(ctx context.Context, uids []int64) (map[int64]domain.FollowStatistic, error) {
	pipeline := c.cmd.Pipeline()
	for _, uid := range uids {
		pipeline.HGetAll(ctx, c.statisticKey(uid))
	}
	cmds, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.FollowStatistic, len(cmds))
	for i, cmd := range cmds {
		data := cmd.(*redis.MapStringStringCmd).Val()
		if len(data) == 0 {
			continue
		}
		followers, _ := strconv.ParseInt(data[fieldFollowerCount], 10, 64)
		following, _ := strconv.ParseInt(data[fieldFolloweeCount], 10, 64)
		res[uids[i]] = domain.FollowStatistic{
			Followers: followers,
			Following: following,
		}
	}
	return res, nil
}

func (c *RedisFollowCache) SetStatistic(ctx context.Context, statistic domain.FollowStatistic) error {
	pipeline := c.cmd.Pipeline()
	pipeline.HMSet(ctx, c.statisticKey(statistic.Uid), fieldFollowerCount, statistic.Followers, fieldFolloweeCount, statistic.Following)
	pipeline.Expire(ctx, c.statisticKey(statistic.Uid), c.exp)
	_, err := pipeline.Exec(ctx)
	return err
}

func (c *RedisFollowCache) SetStatisticBatch(ctx context.Context, stats []domain.FollowStatistic) error {
	pipeline := c.cmd.Pipeline()
	for _, statistic := range stats {
		pipeline.HMSet(ctx, c.statisticKey(statistic.Uid), fieldFollowerCount, statistic.Followers, fieldFolloweeCount, statistic.Following)
		pipeline.Expire(ctx, c.statisticKey(statistic.Uid), c.exp)
	}
	_, err := pipeline.Exec(ctx)
	return err
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
