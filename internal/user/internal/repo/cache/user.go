package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, uid int64, user domain.User) error
	SetBatch(ctx context.Context, users []domain.User) error
	SetByAccount(ctx context.Context, account string, user domain.User) error
	GetByAccount(ctx context.Context, account string) (domain.User, error)
	GetByIds(ctx context.Context, uids []int64) (map[int64]domain.User, error)
	Delete(ctx context.Context, uid int64) error
	DeleteByAccount(ctx context.Context, account string) error
}

var _ UserCache = (*RedisUserCache)(nil)

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 5,
	}
}

func (cache *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

// Get 从缓存中获取用户信息
// 只要error为nil，就认为缓存命中
func (cache *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := cache.key(uid)
	val, err := cache.client.GetEx(ctx, key, cache.expiration).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err

}

func (cache *RedisUserCache) accountKey(account string) string {
	return fmt.Sprintf("user:account:%s", account)
}

func (cache *RedisUserCache) GetByAccount(ctx context.Context, account string) (domain.User, error) {
	val, err := cache.client.GetEx(ctx, cache.accountKey(account), cache.expiration).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) SetByAccount(ctx context.Context, account string, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.accountKey(account)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) DeleteByAccount(ctx context.Context, account string) error {
	key := cache.accountKey(account)
	return cache.client.Del(ctx, key).Err()
}

func (cache *RedisUserCache) GetByIds(ctx context.Context, uids []int64) (map[int64]domain.User, error) {
	keys := make([]string, len(uids))
	for i, uid := range uids {
		keys[i] = cache.key(uid)
	}
	vals, err := cache.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.User)
	for i, val := range vals {
		if val == nil || val == "" {
			continue
		}
		var u domain.User
		err = json.Unmarshal([]byte(val.(string)), &u)
		if err != nil {
			return nil, err
		}
		res[uids[i]] = u
	}
	return res, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, uid int64, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(uid)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) SetBatch(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}
	pipe := cache.client.Pipeline()
	for _, user := range users {
		val, err := json.Marshal(user)
		if err != nil {
			return err
		}
		key := cache.key(user.Id)
		pipe.Set(ctx, key, val, cache.expiration)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (cache *RedisUserCache) Delete(ctx context.Context, uid int64) error {
	key := cache.key(uid)
	return cache.client.Del(ctx, key).Err()
}
