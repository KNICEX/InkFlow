package middleware

import (
	ijwt "github.com/KNICEX/InkFlow/internal/web/jwt"
	"github.com/KNICEX/InkFlow/internal/web/result"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication interface {
	CheckLogin() gin.HandlerFunc
}

type LoginJWTBuilder struct {
	ijwt.Handler
	logger logx.Logger
}

func NewLoginJWTBuilder(handler ijwt.Handler, l logx.Logger) *LoginJWTBuilder {
	return &LoginJWTBuilder{
		Handler: handler,
		logger:  l,
	}
}

func (b *LoginJWTBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := b.Handler.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return ijwt.UserClaimsKey, nil
		})
		if err != nil || token == nil || !token.Valid {
			result.InvalidToken(ctx)
			ctx.Abort()
			return
		}

		ok, err := b.Handler.CheckSession(ctx, claims.Ssid)
		if err != nil {
			b.logger.WithCtx(ctx).Error("login_jwt CheckSession", logx.Error(err))
			result.InternalError(ctx)
			ctx.Abort()
			return
		}
		ctx.Set(ijwt.UserClaimsCtxKey, claims)

		if !ok {
			result.InvalidToken(ctx)
			ctx.Abort()
		}

	}
}
