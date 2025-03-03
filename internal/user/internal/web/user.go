package web

import (
	"errors"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/code"
	"github.com/KNICEX/InkFlow/internal/web/handler"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

const (
	loginBiz    = "login"
	resetPwdBiz = "reset_pwd"
)

var _ handler.Handler = (*Handler)(nil)

type Handler struct {
	svc          service.UserService
	emailCodeSvc code.Service
	phoneReg     *regexp.Regexp
	emailReg     *regexp.Regexp
	l            logx.Logger
	auth         middleware.Authentication
	jwt.Handler
}

func NewUserHandler(svc service.UserService,
	emailCodeSvc code.Service,
	jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) *Handler {
	return &Handler{
		svc:          svc,
		emailCodeSvc: emailCodeSvc,
		phoneReg:     regexp.MustCompile(`^1[3456789]\d{9}$`),
		emailReg:     regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`),
		Handler:      jwtHandler,
		l:            log,
		auth:         auth,
	}
}

func (u *Handler) RegisterRoutes(server *gin.RouterGroup) {

	// 登录验证码
	server.POST("/verify/send/login", ginx.WrapBody(u.l, u.SendLoginCode))

	loginGroup := server.Group("/login")
	{
		// 邮箱验证码登录(自动注册)
		loginGroup.POST("/email", ginx.WrapBody(u.l, u.LoginEmail))
		// 手机验证码登录(自动注册)
		//loginGroup.POST("/sms", ginx.WrapBody(u.l, u.LoginSMS))
		// 邮箱密码登录
		loginGroup.POST("/pwd/email", ginx.WrapBody(u.l, u.LoginEmailPwd))
		// 手机密码登录
		loginGroup.POST("/pwd/phone", ginx.WrapBody(u.l, u.LoginPhonePwd))
		// 用户名密码登录
		loginGroup.POST("/pwd/account_name", ginx.WrapBody(u.l, u.LoginPhonePwd))
	}

	// 刷新token
	server.POST("/refresh_token", u.RefreshToken)

	// 需要登录
	checkGroup := server.Group("")
	checkGroup.Use(u.auth.CheckLogin())
	checkGroup.Use()
	{
		checkGroup.GET("/logout", ginx.Wrap(u.l, u.Logout))

		checkGroup.GET("/profile", ginx.Wrap(u.l, u.Profile))
		checkGroup.PUT("/profile", ginx.WrapBody(u.l, u.EditProfile))
		// 修改账号名
		checkGroup.PUT("/account_name", ginx.WrapBody(u.l, u.EditAccountName))

		// 发送重置密码验证码
		//checkGroup.POST("/verify/send/reset/sms", ginx.Wrap(u.l, u.SendResetPwdCodeSms))
		checkGroup.POST("/verify/send/reset/email", ginx.Wrap(u.l, u.SendResetPwdCodeEmail))

		// 重置密码
		checkGroup.POST("/pwd/reset/email", ginx.WrapBody(u.l, u.ResetPwdByEmailCode))
		//checkGroup.POST("/pwd/reset/sms", ginx.WrapBody(u.l, u.ResetPwdBySmsCode))
		checkGroup.POST("/pwd/reset/old", ginx.WrapBody(u.l, u.ChangePwd))
	}
}

func (u *Handler) sendCodeWithSvc(ctx *gin.Context, biz, recipient string, svc code.Service) (ginx.Result, error) {
	err := svc.Send(ctx, biz, recipient)
	switch {
	case err == nil:
		return ginx.SuccessWithMsg("验证码发送成功"), nil
	case errors.Is(err, code.ErrCodeSendTooMany):
		return ginx.BizError("验证码发送太频繁"), err
	default:
		return ginx.InternalError(), err
	}
}

func (u *Handler) sendCode(ctx *gin.Context, biz string, req SendCodeReq) (ginx.Result, error) {
	if req.Email != "" && u.emailReg.MatchString(req.Email) {
		return u.sendCodeWithSvc(ctx, biz, req.Email, u.emailCodeSvc)
	}
	//if req.Phone != "" && u.phoneReg.MatchString(req.Phone) {
	//	return u.sendCodeWithSvc(ctx, biz, req.Phone, u.smsCodeSvc)
	//}
	return ginx.InvalidParam(), nil
}

func (u *Handler) verifyCode(ctx *gin.Context, biz, recipient, verifyCode string, svc code.Service) (ginx.Result, error) {
	ok, err := svc.Verify(ctx, biz, recipient, verifyCode)
	switch {
	case err != nil && !errors.Is(err, code.ErrCodeVerifyLimit):
		return ginx.InternalError(), err
	case !ok:
		return ginx.BizError("验证码错误"), nil
	default:
		return ginx.Success(), nil
	}
}

//func (u *Handler) SendResetPwdCodeSms(ctx *gin.Context) (ginx.Result, error) {
//	uc := ijwt.MustGetUserClaims(ctx)
//	user, err := u.svc.Profile(ctx, uc.UserId)
//	if err != nil {
//		return ginx.InternalError(), err
//	}
//	if user.Phone == "" {
//		return ginx.BizError("用户未绑定手机号"), nil
//	}
//	return u.sendCodeWithSvc(ctx, resetPwdBiz, user.Phone, u.smsCodeSvc)
//}

func (u *Handler) SendResetPwdCodeEmail(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	user, err := u.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	if user.Email == "" {
		return ginx.BizError("用户未绑定邮箱"), nil
	}
	return u.sendCodeWithSvc(ctx, resetPwdBiz, user.Email, u.emailCodeSvc)
}

func (u *Handler) SendLoginCode(ctx *gin.Context, req SendCodeReq) (ginx.Result, error) {
	return u.sendCode(ctx, loginBiz, req)
}

// LoginSMS 手机验证码登录 登录或自动创建用户
//func (u *Handler) LoginSMS(ctx *gin.Context, req LoginSmsReq) (ginx.Result, error) {
//	ok, err := u.smsCodeSvc.Verify(ctx, loginBiz, req.Phone, req.Code)
//	if err != nil && !errors.Is(err, code.ErrCodeVerifyLimit) {
//		return ginx.InternalError(), err
//	}
//	if !ok {
//		return ginx.BizError("验证码错误"), nil
//	}
//
//	user, err := u.svc.FindOrCreateByPhone(ctx, req.Phone)
//	if err != nil {
//		return ginx.InternalError(), err
//	}
//
//	if err = u.SetLoginToken(ctx, user); err != nil {
//		return ginx.InternalError(), err
//	}
//
//	return ginx.SuccessWithMsg("登录成功"), nil
//}

// LoginPhonePwd 手机号密码登录
func (u *Handler) LoginPhonePwd(ctx *gin.Context, req LoginPhonePwdReq) (ginx.Result, error) {
	user, err := u.svc.LoginPhonePwd(ctx, req.Phone, req.Password)
	if errors.Is(err, service.ErrInvalidAccountOrPwd) {
		return ginx.BizError("用户或密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		return ginx.InternalError(), err
	}

	return ginx.SuccessWithMsg("登录成功"), nil
}

func (u *Handler) LoginAccountNamePwd(ctx *gin.Context, req LoginAccountNamePwdReq) (ginx.Result, error) {
	user, err := u.svc.LoginAccountNamePwd(ctx, req.AccountName, req.Password)
	if errors.Is(err, service.ErrInvalidAccountOrPwd) {
		return ginx.BizError("账号名或密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		return ginx.InternalError(), err
	}

	return ginx.SuccessWithMsg("登录成功"), nil
}

// LoginEmail 邮箱验证码登录(自动注册)
func (u *Handler) LoginEmail(ctx *gin.Context, req LoginEmailReq) (ginx.Result, error) {
	var user domain.User
	ok, err := u.emailCodeSvc.Verify(ctx, loginBiz, req.Email, req.Code)
	if err != nil && !errors.Is(err, code.ErrCodeVerifyLimit) {
		return ginx.InternalError(), err
	}
	if !ok {
		return ginx.BizError("验证码错误"), nil
	}

	user, err = u.svc.FindOrCreateByEmail(ctx, req.Email)
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登录成功"), nil
}

// LoginEmailPwd 邮箱密码登录
func (u *Handler) LoginEmailPwd(ctx *gin.Context, req LoginPhonePwdReq) (ginx.Result, error) {
	user, err := u.svc.LoginEmailPwd(ctx, req.Phone, req.Password)
	if errors.Is(err, service.ErrInvalidAccountOrPwd) {
		return ginx.BizError("邮箱或密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		return ginx.InternalError(), err
	}

	return ginx.SuccessWithMsg("登录成功"), nil
}

func (u *Handler) Logout(ctx *gin.Context) (ginx.Result, error) {
	if err := u.ClearToken(ctx); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登出成功"), nil
}

func (u *Handler) Profile(ctx *gin.Context) (ginx.Result, error) {
	type Profile struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		AccountName string `json:"accountName"`
		Username    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		AboutMe     string `json:"aboutMe"`
	}

	uc := jwt.MustGetUserClaims(ctx)

	user, err := u.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	var birthday string
	if !user.Birthday.IsZero() {
		birthday = user.Birthday.Format(time.DateOnly)
	}
	return ginx.SuccessWithData(Profile{
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		AccountName: user.Account,
		Birthday:    birthday,
		AboutMe:     user.AboutMe,
	}), nil
}

func (u *Handler) EditAccountName(ctx *gin.Context, req EditAccountNameReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	err := u.svc.UpdateAccountName(ctx, uc.UserId, req.AccountName)
	if errors.Is(err, service.ErrUserDuplicate) {
		return ginx.BizError("账号名已存在"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("修改成功"), nil
}

func (u *Handler) EditProfile(ctx *gin.Context, req EditProfileReq) (ginx.Result, error) {
	var birthday time.Time
	var err error
	if req.Birthday != "" {
		birthday, err = time.Parse(time.DateOnly, req.Birthday)
		if err != nil {
			return ginx.InvalidParamWithMsg("生日格式错误"), nil
		}
	}
	uc := jwt.MustGetUserClaims(ctx)
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.UserId,
		Username: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("修改成功"), nil
}

func (u *Handler) RefreshToken(ctx *gin.Context) {
	type Req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	rc, err := jwt.ParseClaims(req.RefreshToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if rc.ExpiresAt.Time.Sub(time.Now()) < time.Hour*24*7 {
		// refresh token 也快过期了，一起刷新
		if err = u.ClearToken(ctx); err != nil {
			u.l.WithCtx(ctx).Error("RefreshToken 一并刷新refreshToken失败", logx.Error(err))
		}
		err = u.SetLoginToken(ctx, rc.UserId)
	} else {
		// 只刷新短token
		err = u.SetJwtToken(ctx, rc.UserId, rc.Ssid)
	}
	if err != nil {
		u.l.WithCtx(ctx).Error("RefreshToken", logx.Error(err))
		ctx.JSON(http.StatusOK, ginx.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, ginx.Success())
}

func (u *Handler) resetPwdByCode(ctx *gin.Context, verify func(user domain.User) (bool, error)) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	user, err := u.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	ok, err := verify(user)
	if err != nil {
		return ginx.InternalError(), err
	}
	if !ok {
		return ginx.BizError("验证码错误"), nil
	}
	if err = u.ClearToken(ctx); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("重置密码成功, 请重新登录"), nil
}

//func (u *Handler) ResetPwdBySmsCode(ctx *gin.Context, req SmsResetPwdReq) (ginx.Result, error) {
//	return u.resetPwdByCode(ctx, func(user domain.User) (bool, error) {
//		if user.Phone == "" {
//			return false, errors.New("用户未绑定手机号")
//		}
//		return u.smsCodeSvc.Verify(ctx, resetPwdBiz, user.Phone, req.Code)
//	})
//}

func (u *Handler) ResetPwdByEmailCode(ctx *gin.Context, req EmailResetPwdReq) (ginx.Result, error) {
	return u.resetPwdByCode(ctx, func(user domain.User) (bool, error) {
		if user.Email == "" {
			return false, errors.New("用户未绑定邮箱")
		}
		return u.emailCodeSvc.Verify(ctx, resetPwdBiz, user.Email, req.Code)
	})
}

func (u *Handler) ChangePwd(ctx *gin.Context, req ChangePwdReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	err := u.svc.ChangePwd(ctx, uc.UserId, req.OldPassword, req.NewPassword)
	if errors.Is(err, service.ErrInvalidAccountOrPwd) {
		return ginx.BizError("原密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("密码修改成功"), nil
}
