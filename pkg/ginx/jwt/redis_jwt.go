package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	RefreshClaimsKey = []byte("refresh-key")
	UserClaimsKey    = []byte("user-key")
)

const (
	UserClaimsCtxKey      = "_uid"
	TokenHeaderKey        = "x-access-token"
	RefreshTokenHeaderKey = "x-refresh-token"
)

type RedisHandler struct {
	client        redis.Cmdable
	refreshExpire time.Duration
}

func NewRedisHandler(client redis.Cmdable) Handler {
	return &RedisHandler{
		client: client,
		// TODO 考虑从配置文件中读取
		refreshExpire: time.Hour * 24 * 30,
	}
}

func (h *RedisHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	if err := h.setRefreshToken(ctx, uid, ssid); err != nil {
		return err
	}
	return h.SetJwtToken(ctx, uid, ssid)
}

func (h *RedisHandler) SetJwtToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		UserId: uid,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(UserClaimsKey)
	if err != nil {
		return err
	}
	ctx.Header(TokenHeaderKey, tokenStr)
	return nil
}

func (h *RedisHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		UserId: uid,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.refreshExpire)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	tokenStr, err := token.SignedString(RefreshClaimsKey)
	if err != nil {
		return err
	}
	ctx.Header(RefreshTokenHeaderKey, tokenStr)
	return nil
}

func (h *RedisHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (h *RedisHandler) CheckSession(ctx *gin.Context, ssid string) (bool, error) {
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil {
		return false, err
	}
	// redis保存已失效的session
	if cnt > 0 {
		return false, nil
	}
	return true, nil
}
func (h *RedisHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header(TokenHeaderKey, "")
	ctx.Header(RefreshTokenHeaderKey, "")
	uc := ctx.MustGet(UserClaimsCtxKey).(UserClaims)
	// 标记session失效
	return h.client.Set(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid), "1", h.refreshExpire).Err()
}
