package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type InteractiveCache interface {
	IncrViewCnt(ctx context.Context, biz string, bizId int64) error

	IncrViewCntBatch(ctx context.Context, biz string, bizIds []int64) error
	IncrLikeCnt(ctx context.Context, biz string, bizId int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizId int64) error
	IncrFavoriteCnt(ctx context.Context, biz string, bizId int64) error
	DecrFavoriteCnt(ctx context.Context, biz string, bizId int64) error

	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	GetBatch(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
	SetBatch(ctx context.Context, intrs []domain.Interactive) error
}

var (
	//go:embed lua/interactive_incr_cnt.lua
	luaIncrCnt string

	ErrKeyNotFound = redis.Nil
)

const (
	fieldReadCnt     = "read_cnt"
	fieldFavoriteCnt = "favorite_cnt"
	fieldLikeCnt     = "like_cnt"
)

type RedisInteractiveCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisInteractiveCache(cmd redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		cmd:        cmd,
		expiration: time.Minute * 3,
	}
}
func (cache *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (cache *RedisInteractiveCache) IncrViewCnt(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldReadCnt, 1).Err()
}
func (cache *RedisInteractiveCache) IncrViewCntBatch(ctx context.Context, biz string, bizIds []int64) error {
	pipeline := cache.cmd.Pipeline()
	for _, bizId := range bizIds {
		pipeline.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldReadCnt, 1)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func (cache *RedisInteractiveCache) IncrLikeCnt(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, 1).Err()
}

func (cache *RedisInteractiveCache) DecrLikeCnt(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, -1).Err()
}

func (cache *RedisInteractiveCache) IncrFavoriteCnt(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldFavoriteCnt, 1).Err()
}

func (cache *RedisInteractiveCache) DecrFavoriteCnt(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldFavoriteCnt, -1).Err()
}

func (cache *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := cache.key(biz, bizId)
	pipeline := cache.cmd.Pipeline()
	pipeline.HMSet(ctx, key,
		fieldReadCnt, intr.ViewCnt,
		fieldLikeCnt, intr.LikeCnt,
		fieldFavoriteCnt, intr.CollectCnt)

	pipeline.Expire(ctx, key, cache.expiration)
	_, err := pipeline.Exec(ctx)
	return err
}

func (cache *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	data, err := cache.cmd.HGetAll(ctx, cache.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(data) == 0 {
		return domain.Interactive{}, ErrKeyNotFound
	}

	collectCnt, _ := strconv.ParseInt(data[fieldFavoriteCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)
	return domain.Interactive{
		Biz:        biz,
		BizId:      bizId,
		ViewCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, nil
}

func (cache *RedisInteractiveCache) GetBatch(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	pipeline := cache.cmd.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(bizIds))
	for i, bizId := range bizIds {
		cmds[i] = pipeline.HGetAll(ctx, cache.key(biz, bizId))
	}
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(bizIds))
	for i, cmd := range cmds {
		data := cmd.Val()
		if len(data) == 0 {
			continue
		}
		bizId := bizIds[i]
		collectCnt, _ := strconv.ParseInt(data[fieldFavoriteCnt], 10, 64)
		likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
		readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)
		res[bizId] = domain.Interactive{
			Biz:        biz,
			BizId:      bizId,
			ViewCnt:    readCnt,
			LikeCnt:    likeCnt,
			CollectCnt: collectCnt,
		}
	}
	return res, nil
}

func (cache *RedisInteractiveCache) SetBatch(ctx context.Context, intrs []domain.Interactive) error {
	pipeline := cache.cmd.Pipeline()
	for _, intr := range intrs {
		key := cache.key(intr.Biz, intr.BizId)
		pipeline.HMSet(ctx, key,
			fieldReadCnt, intr.ViewCnt,
			fieldLikeCnt, intr.LikeCnt,
			fieldFavoriteCnt, intr.CollectCnt)
		pipeline.Expire(ctx, key, cache.expiration)
	}
	_, err := pipeline.Exec(ctx)
	return err
}
