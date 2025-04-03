package cache

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type CommentCache interface {
	IncrLike(ctx context.Context, id int64) error
	DecrLike(ctx context.Context, id int64) error
	IncrReply(ctx context.Context, id int64) error
	DecrReply(ctx context.Context, id int64) error

	IncrBizReply(ctx context.Context, biz string, bizId int64) error
	DecrBizReply(ctx context.Context, biz string, bizId int64) error

	SetBizReplyCount(ctx context.Context, biz string, counts map[int64]int64) error
	BizReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error)
}

var (
	//go:embed lua/incr_exist.lua
	incrExistLua string
)

type RedisCommentCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisCommentCache(cmd redis.Cmdable) CommentCache {
	return &RedisCommentCache{
		cmd:        cmd,
		expiration: time.Hour,
	}
}

func (cache *RedisCommentCache) bizReplyKey(biz string, bizId int64) string {
	return "comment:reply:" + biz + ":" + strconv.FormatInt(bizId, 10)
}

func (cache *RedisCommentCache) IncrLike(ctx context.Context, id int64) error {
	return nil
}

func (cache *RedisCommentCache) DecrLike(ctx context.Context, id int64) error {
	return nil
}

func (cache *RedisCommentCache) IncrReply(ctx context.Context, id int64) error {
	return nil
}

func (cache *RedisCommentCache) DecrReply(ctx context.Context, id int64) error {
	return nil
}

func (cache *RedisCommentCache) IncrBizReply(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, incrExistLua, []string{cache.bizReplyKey(biz, bizId)}, 1, cache.expiration.Seconds()).Err()
}

func (cache *RedisCommentCache) DecrBizReply(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, incrExistLua, []string{cache.bizReplyKey(biz, bizId)}, -1, cache.expiration.Seconds()).Err()
}

func (cache *RedisCommentCache) SetBizReplyCount(ctx context.Context, biz string, counts map[int64]int64) error {
	if len(counts) == 0 {
		return nil
	}
	pipe := cache.cmd.Pipeline()
	for bizId, count := range counts {
		pipe.Set(ctx, cache.bizReplyKey(biz, bizId), count, cache.expiration)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (cache *RedisCommentCache) BizReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error) {
	res, err := cache.cmd.MGet(ctx, lo.Map(bizIds, func(item int64, index int) string {
		return cache.bizReplyKey(biz, item)
	})...).Result()
	if err != nil {
		return nil, err
	}

	replyCountMap := make(map[int64]int64, len(res))
	for i, item := range res {
		if item == nil || item == "" {
			continue
		}
		replyCount, err := strconv.ParseInt(item.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		replyCountMap[bizIds[i]] = replyCount
	}
	return replyCountMap, nil
}
