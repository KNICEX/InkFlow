package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type CommentCache interface {
	IncrLike(ctx context.Context, id int64) error
	DecrLike(ctx context.Context, id int64) error
	IncrReply(ctx context.Context, id int64) error
	DecrReply(ctx context.Context, id int64) error
}

type RedisCommentCache struct {
	cmd redis.Cmdable
}

func NewRedisCommentCache(cmd redis.Cmdable) CommentCache {
	return &RedisCommentCache{
		cmd: cmd,
	}
}

func (r *RedisCommentCache) IncrLike(ctx context.Context, id int64) error {
	return nil
}

func (r *RedisCommentCache) DecrLike(ctx context.Context, id int64) error {
	return nil
}

func (r *RedisCommentCache) IncrReply(ctx context.Context, id int64) error {
	return nil
}

func (r *RedisCommentCache) DecrReply(ctx context.Context, id int64) error {
	return nil
}
