package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var (
	ErrCodeSendTooMany = errors.New("发送太频繁")
	ErrCodeVerifyLimit = errors.New("超出验证次数")
	ErrUnknownForCode  = errors.New("未知错误")
)

var _ CodeCache = (*RedisCodeCache)(nil)

type CodeCache interface {
	Set(ctx context.Context, biz, recipient, code string, effectiveTime time.Duration, resendInterval time.Duration, maxRetry int) error
	Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, recipient, code string, effectiveTime time.Duration, resendInterval time.Duration, maxRetry int) error {
	if effectiveTime == 0 {
		effectiveTime = time.Minute * 5
	}

	effectiveEx := effectiveTime.Seconds()
	resendEx := resendInterval.Seconds()
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, recipient)}, code, effectiveEx, resendEx, maxRetry).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//	成功
		return nil
	case -1:
		//	发送太频繁
		//zap.L().Warn("验证码太频繁", zap.String("biz", biz),
		//	zap.String("recipient", recipient))
		return ErrCodeSendTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, recipient)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		//	验证成功
		return true, nil
	case -1:
		//	验证次数用完/验证码过期/验证码被使用过/验证码不存在
		return false, ErrCodeVerifyLimit
	case -2:
		//	规定次数内输入错误
		return false, nil
	default:
		// 应该不会出现
		return false, ErrUnknownForCode
	}

}

func (c *RedisCodeCache) key(biz, recipient string) string {
	return fmt.Sprintf("recipient_code:%s:%s", biz, recipient)
}
