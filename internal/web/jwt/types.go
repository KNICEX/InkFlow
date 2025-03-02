package jwt

import (
	"github.com/KNICEX/InkFlow/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetJwtToken(ctx *gin.Context, user domain.User, ssid string) error
	SetLoginToken(ctx *gin.Context, user domain.User) error
	CheckSession(ctx *gin.Context, ssid string) (bool, error)

	ClearToken(ctx *gin.Context) error
}

func MustGetUserClaims(ctx *gin.Context) UserClaims {
	return ctx.MustGet(UserClaimsCtxKey).(UserClaims)
}

func GetUserClaims(ctx *gin.Context) (UserClaims, bool) {
	claims, ok := ctx.Get(UserClaimsCtxKey)
	if !ok {
		return UserClaims{}, false
	}
	return claims.(UserClaims), true
}

type UserClaims struct {
	jwt.RegisteredClaims
	Ssid   string `json:"ssid"`
	UserId int64  `json:"userId"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Ssid   string `json:"ssid"`
	UserId int64  `json:"userId"`
}
