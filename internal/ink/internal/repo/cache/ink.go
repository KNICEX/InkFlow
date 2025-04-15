package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"time"
)

var (
	ErrKeyNotFound = redis.Nil
)

type InkCache interface {
	Get(ctx context.Context, id int64) (domain.Ink, error)
	Set(ctx context.Context, ink domain.Ink) error
	SetBatch(ctx context.Context, inks []domain.Ink) error
	SetFirstPage(ctx context.Context, authorId int64, inks []domain.Ink) error
	GetFirstPage(ctx context.Context, authorId int64) ([]domain.Ink, error)
	DelFirstPage(ctx context.Context, authorId int64) error
	SetPage(ctx context.Context, authorId int64, offset, limit int, inks []domain.Ink) error
	GetPage(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error)
	DelPage(ctx context.Context, authorId int64, offset, limit int) error
	GetByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error)
	Del(ctx context.Context, id int64) error
}

type RedisInkCache struct {
	cmd    redis.Cmdable
	expire time.Duration
}

func NewRedisInkCache(cmd redis.Cmdable) InkCache {
	return &RedisInkCache{
		cmd:    cmd,
		expire: time.Minute * 10,
	}
}

func (cache *RedisInkCache) key(id int64) string {
	return fmt.Sprintf("ink:detail:%d", id)
}

func (cache *RedisInkCache) firstPageKey(authorId int64) string {
	return fmt.Sprintf("ink:first-page:%d", authorId)
}

func (cache *RedisInkCache) pageKey(authorId int64, offset, limit int) string {
	return fmt.Sprintf("ink:page:%d:%d:%d", authorId, offset, limit)
}

func (cache *RedisInkCache) Get(ctx context.Context, id int64) (domain.Ink, error) {
	ink := domain.Ink{}
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return ink, err
	}
	err = json.Unmarshal(val, &ink)
	if err != nil {
		return ink, err
	}
	return ink, nil
}

func (cache *RedisInkCache) Set(ctx context.Context, ink domain.Ink) error {
	key := cache.key(ink.Id)
	val, err := json.Marshal(ink)
	if err != nil {
		return err
	}
	err = cache.cmd.Set(ctx, key, val, cache.expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) SetBatch(ctx context.Context, inks []domain.Ink) error {
	if len(inks) == 0 {
		return nil
	}
	pipe := cache.cmd.Pipeline()
	for _, ink := range inks {
		key := cache.key(ink.Id)
		val, err := json.Marshal(ink)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, val, cache.expire)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) SetFirstPage(ctx context.Context, authorId int64, inks []domain.Ink) error {
	for i, _ := range inks {
		inks[i].ContentHtml = inks[i].Abstract()
	}
	val, err := json.Marshal(inks)
	if err != nil {
		return err
	}
	err = cache.cmd.Set(ctx, cache.firstPageKey(authorId), val, cache.expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) GetFirstPage(ctx context.Context, authorId int64) ([]domain.Ink, error) {
	val, err := cache.cmd.Get(ctx, cache.firstPageKey(authorId)).Bytes()
	if err != nil {
		return nil, err
	}
	var inks []domain.Ink
	err = json.Unmarshal(val, &inks)
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (cache *RedisInkCache) SetPage(ctx context.Context, authorId int64, offset, limit int, inks []domain.Ink) error {
	// 这里不会保存文章全量内容， 只截取前面的部分
	for i, _ := range inks {
		inks[i].ContentHtml = inks[i].Abstract()
	}
	val, err := json.Marshal(inks)
	if err != nil {
		return err
	}
	err = cache.cmd.Set(ctx, cache.pageKey(authorId, offset, limit), val, cache.expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) GetPage(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	val, err := cache.cmd.Get(ctx, cache.pageKey(authorId, offset, limit)).Bytes()
	if err != nil {
		return nil, err
	}
	var inks []domain.Ink
	err = json.Unmarshal(val, &inks)
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (cache *RedisInkCache) DelPage(ctx context.Context, authorId int64, offset, limit int) error {
	key := cache.pageKey(authorId, offset, limit)
	err := cache.cmd.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) GetByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error) {
	keys := lo.Map(ids, func(item int64, index int) string {
		return cache.key(item)
	})
	vals, err := cache.cmd.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	inkMap := make(map[int64]domain.Ink, len(ids))
	for i, val := range vals {
		if val == nil || val == "" {
			continue
		}
		var ink domain.Ink
		err = json.Unmarshal([]byte(val.(string)), &ink)
		if err != nil {
			return nil, err
		}
		inkMap[ids[i]] = ink
	}
	if len(inkMap) == 0 {
		return nil, ErrKeyNotFound
	}
	return inkMap, nil
}

func (cache *RedisInkCache) DelFirstPage(ctx context.Context, authorId int64) error {
	key := cache.firstPageKey(authorId)
	err := cache.cmd.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisInkCache) Del(ctx context.Context, id int64) error {
	key := cache.key(id)
	err := cache.cmd.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
