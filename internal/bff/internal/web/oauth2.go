package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	ijwt "github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/uuidx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const oauth2CookieName = "jwt-state"

type OAuth2Handler[T any] struct {
	svc          user.OAuth2Service[T]
	stateKey     string
	callBackPath string
	getDomain    func(ctx context.Context, t T) (user.User, error)
	l            logx.Logger
	ijwt.Handler
}

func (o *OAuth2Handler[T]) AuthUrl(ctx *gin.Context) (ginx.Result, error) {
	state := uuidx.NewShort()
	url, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		o.l.WithCtx(ctx).Error("oauth get auth url failed", logx.Error(err))
		return ginx.InternalError(), err
	}
	if err = o.setStatCookie(ctx, state); err != nil {
		o.l.WithCtx(ctx).Error("oauth set state cookie failed", logx.Error(err))
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(url), nil
}

func (o *OAuth2Handler[T]) setStatCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString([]byte(o.stateKey))
	if err != nil {
		return err
	}
	// TODO 解决SameSite 问题或者换实现方式
	cookie := &http.Cookie{
		Name:    oauth2CookieName,
		Value:   tokenStr,
		Expires: time.Now().Add(time.Minute * 10),
		Path:    o.callBackPath,
		// 只能通过https访问
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(ctx.Writer, cookie)
	return nil
}

func (o *OAuth2Handler[T]) verifyState(ctx *gin.Context) error {
	//state := ctx.Query("state")
	//// 验证state
	//cookieState, err := ctx.Cookie(oauth2CookieName)
	//if err != nil {
	//	return fmt.Errorf("get cookie state failed, %w", err)
	//}
	//var sc StateClaims
	//token, err := jwt.ParseWithClaims(cookieState, &sc, func(token *jwt.Token) (any, error) {
	//	return []byte(o.stateKey), nil
	//})
	//if err != nil || !token.Valid {
	//	return fmt.Errorf("token has expired, %w", err)
	//}
	//if sc.State != state {
	//	return errors.New("state not match")
	//}
	return nil

}

func (o *OAuth2Handler[T]) Callback(ctx *gin.Context, callback Oauth2Callback) (ginx.Result, error) {
	if callback.Code == "" {
		return ginx.InvalidParam(), nil
	}

	err := o.verifyState(ctx)
	if err != nil {
		o.l.WithCtx(ctx).Error("oauth2 verify state failed", logx.Error(err))
		return ginx.InvalidParamWithMsg("登录失败"), err
	}

	info, err := o.svc.VerifyCode(ctx, callback.Code)
	if err != nil {
		o.l.WithCtx(ctx).Error("oauth2 verify code failed", logx.Error(err))
		return ginx.InternalError(), err
	}

	// 通过wechat/GitHub Info 获取domain.User
	u, err := o.getDomain(ctx, info)

	if err != nil {
		o.l.WithCtx(ctx).Error("oauth2 get domain user failed", logx.Error(err))
		return ginx.InternalError(), err
	}
	if err = o.SetLoginToken(ctx, u.Id); err != nil {
		o.l.WithCtx(ctx).Error("oauth2 set login token failed", logx.Error(err))
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登录成功"), nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}

type Config struct {
	Secure bool
}
