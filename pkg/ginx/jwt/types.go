package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetJwtToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	CheckSession(ctx *gin.Context, ssid string) (bool, error)

	ClearToken(ctx *gin.Context) error
}

func MustGetUserClaims(ctx *gin.Context) UserClaims {
	return ctx.MustGet(UserClaimsCtxKey).(UserClaims)
}

func ParseClaims(tokenStr string) (UserClaims, error) {
	res := UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &res, func(token *jwt.Token) (any, error) {
		return UserClaimsKey, nil
	})
	if err != nil {
		return res, err
	}
	if !token.Valid {
		return res, errors.New("token invalid")
	}
	return res, nil
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
