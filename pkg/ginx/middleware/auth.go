package middleware

import (
	ijwt "github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Authentication interface {
	CheckLogin() gin.HandlerFunc
}

type JwtLoginBuilder struct {
	ijwt.Handler
	l logx.Logger
}

func (b *JwtLoginBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := b.Handler.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return ijwt.UserClaimsKey, nil
		})
		if err != nil || token == nil || !token.Valid {

			ctx.Abort()
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
