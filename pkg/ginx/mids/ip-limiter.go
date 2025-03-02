package mids

import (
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type IpLimiterBuilder struct {
	prefix  string
	limiter ratelimit.KeyLimiter
}

func NewIpLimiterBuilder(limiter ratelimit.KeyLimiter, opts ...Option) *IpLimiterBuilder {
	res := &IpLimiterBuilder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

type Option func(b *IpLimiterBuilder)

func WithPrefix(prefix string) Option {
	return func(b *IpLimiterBuilder) {
		b.prefix = prefix
	}
}

func (b *IpLimiterBuilder) WithPrefix(prefix string) *IpLimiterBuilder {
	b.prefix = prefix
	return b
}

func (b *IpLimiterBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			// 这一步很有意思，就是如果这边出错了
			// 要怎么办？
			// 保守做法：因为借助于 Redis 来做限流，那么 Redis 崩溃了，为了防止系统崩溃，直接限流
			ctx.AbortWithStatus(http.StatusInternalServerError)
			// 激进做法：虽然 Redis 崩溃了，但是这个时候还是要尽量服务正常的用户，所以不限流
			// ctx.Next()
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *IpLimiterBuilder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.limiter.Limited(ctx, key)
}
