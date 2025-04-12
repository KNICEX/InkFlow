package middleware

import (
	ijwt "github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Authentication interface {
	CheckLogin() gin.HandlerFunc     // 检查登录，并且设置用户信息到上下文
	ExtractPayload() gin.HandlerFunc // 提取用户信息到上下文，未登录的用户也可以访问
}

type JwtLoginBuilder struct {
	ijwt.Handler
	l logx.Logger
}

func NewJwtLoginBuilder(handler ijwt.Handler, l logx.Logger) Authentication {
	return &JwtLoginBuilder{
		Handler: handler,
		l:       l,
	}
}

func (b *JwtLoginBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 预检ctx是否已经有用户信息
		if _, ok := ctx.Get(ijwt.UserClaimsCtxKey); ok {
			ctx.Next()
			return
		}

		tokenStr := b.Handler.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return ijwt.UserClaimsKey, nil
		})
		if err != nil || token == nil || !token.Valid {

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检查当前token是否有效
		ok, err := b.Handler.CheckSession(ctx, claims.Ssid)
		if err != nil {
			b.l.WithCtx(ctx).Error("login_jwt CheckSession", logx.Error(err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(ijwt.UserClaimsCtxKey, claims)
	}
}

func (b *JwtLoginBuilder) ExtractPayload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := b.Handler.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return ijwt.UserClaimsKey, nil
		})
		if err != nil || token == nil || !token.Valid {
			ctx.Next()
			return
		}

		// 检查当前token是否有效
		ok, err := b.Handler.CheckSession(ctx, claims.Ssid)
		if err != nil || !ok {
			ctx.Next()
			return
		}
		ctx.Set(ijwt.UserClaimsCtxKey, claims)
	}
}
