package web

import (
	"errors"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/user"
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

var _ ginx.Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc      user.Service
	codeSvc  code.Service
	phoneReg *regexp.Regexp
	emailReg *regexp.Regexp
	l        logx.Logger
	auth     middleware.Authentication
	jwt.Handler
}

func NewUserHandler(svc user.Service,
	codeSvc code.Service,
	jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) *UserHandler {
	return &UserHandler{
		svc:      svc,
		codeSvc:  codeSvc,
		phoneReg: regexp.MustCompile(`^1[3456789]\d{9}$`),
		emailReg: regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`),
		Handler:  jwtHandler,
		l:        log,
		auth:     auth,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/user")
	// 登录验证码
	userGroup.POST("/verify/send/login", ginx.WrapBody(h.l, h.SendLoginCode))

	loginGroup := userGroup.Group("/login")
	{
		// 邮箱验证码登录(自动注册)
		loginGroup.POST("/email", ginx.WrapBody(h.l, h.LoginEmail))
		// 邮箱密码登录
		loginGroup.POST("/pwd/email", ginx.WrapBody(h.l, h.LoginEmailPwd))
		// 账号密码登录
		loginGroup.POST("/pwd/account", ginx.WrapBody(h.l, h.LoginAccountPwd))
	}

	// 刷新token
	userGroup.POST("/refresh_token", h.RefreshToken)

	// 需要登录
	checkGroup := userGroup.Group("")
	checkGroup.Use(h.auth.CheckLogin())
	checkGroup.Use()
	{
		checkGroup.GET("/logout", ginx.Wrap(h.l, h.Logout))

		checkGroup.GET("/profile", ginx.Wrap(h.l, h.Profile))
		checkGroup.PUT("/profile", ginx.WrapBody(h.l, h.EditProfile))
		// 修改账号名
		checkGroup.PUT("/account_name", ginx.WrapBody(h.l, h.EditAccountName))

		// 发送重置密码验证码
		//checkGroup.POST("/verify/send/reset/sms", ginx.Wrap(h.l, h.SendResetPwdCodeSms))
		checkGroup.POST("/verify/send/reset/email", ginx.Wrap(h.l, h.SendResetPwdCodeEmail))

		// 重置密码
		checkGroup.POST("/pwd/reset/email", ginx.WrapBody(h.l, h.ResetPwdByEmailCode))
		//checkGroup.POST("/pwd/reset/sms", ginx.WrapBody(h.l, h.ResetPwdBySmsCode))
		checkGroup.POST("/pwd/reset/old", ginx.WrapBody(h.l, h.ChangePwd))
	}
}

func (h *UserHandler) sendCodeWithSvc(ctx *gin.Context, biz, recipient string) (ginx.Result, error) {
	err := h.codeSvc.Send(ctx, biz, recipient)
	switch {
	case err == nil:
		return ginx.SuccessWithMsg("验证码发送成功"), nil
	case errors.Is(err, code.ErrCodeSendTooMany):
		return ginx.BizError("验证码发送太频繁"), err
	default:
		return ginx.InternalError(), err
	}
}

func (h *UserHandler) sendCode(ctx *gin.Context, biz string, req SendCodeReq) (ginx.Result, error) {
	if req.Email != "" && h.emailReg.MatchString(req.Email) {
		return h.sendCodeWithSvc(ctx, biz, req.Email)
	}
	return ginx.InvalidParam(), nil
}

func (h *UserHandler) verifyCode(ctx *gin.Context, biz, recipient, verifyCode string) (ginx.Result, error) {
	ok, err := h.codeSvc.Verify(ctx, biz, recipient, verifyCode)
	switch {
	case err != nil && !errors.Is(err, code.ErrCodeVerifyLimit):
		return ginx.InternalError(), err
	case !ok:
		return ginx.BizError("验证码错误"), nil
	default:
		return ginx.Success(), nil
	}
}

func (h *UserHandler) SendResetPwdCodeEmail(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	u, err := h.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	if u.Email == "" {
		return ginx.BizError("用户未绑定邮箱"), nil
	}
	return h.sendCodeWithSvc(ctx, resetPwdBiz, u.Email)
}

func (h *UserHandler) SendLoginCode(ctx *gin.Context, req SendCodeReq) (ginx.Result, error) {
	return h.sendCode(ctx, loginBiz, req)
}

// LoginEmail 邮箱验证码登录(自动注册)
func (h *UserHandler) LoginEmail(ctx *gin.Context, req LoginEmailReq) (ginx.Result, error) {
	var u user.User
	ok, err := h.codeSvc.Verify(ctx, loginBiz, req.Email, req.Code)
	if err != nil && !errors.Is(err, code.ErrCodeVerifyLimit) {
		return ginx.InternalError(), err
	}
	if !ok {
		return ginx.BizError("验证码错误"), nil
	}

	u, err = h.svc.FindOrCreateByEmail(ctx, req.Email)
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登录成功"), nil
}

// LoginEmailPwd 邮箱密码登录
func (h *UserHandler) LoginEmailPwd(ctx *gin.Context, req LoginEmailPwdReq) (ginx.Result, error) {
	u, err := h.svc.LoginEmailPwd(ctx, req.Email, req.Password)
	if errors.Is(err, user.ErrInvalidAccountOrPwd) {
		return ginx.BizError("邮箱或密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}

	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		return ginx.InternalError(), err
	}

	return ginx.SuccessWithMsg("登录成功"), nil
}

func (h *UserHandler) Logout(ctx *gin.Context) (ginx.Result, error) {
	if err := h.ClearToken(ctx); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登出成功"), nil
}

func (h *UserHandler) LoginAccountPwd(ctx *gin.Context, req LoginAccountPwdReq) (ginx.Result, error) {
	u, err := h.svc.LoginAccountPwd(ctx, req.Account, req.Password)
	if errors.Is(err, user.ErrInvalidAccountOrPwd) {
		return ginx.BizError("账号名或密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}
	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登录成功"), nil
}

func (h *UserHandler) Profile(ctx *gin.Context) (ginx.Result, error) {
	type Profile struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		AccountName string `json:"accountName"`
		Username    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		AboutMe     string `json:"aboutMe"`
	}

	uc := jwt.MustGetUserClaims(ctx)

	u, err := h.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	var birthday string
	if !u.Birthday.IsZero() {
		birthday = u.Birthday.Format(time.DateOnly)
	}
	return ginx.SuccessWithData(Profile{
		Email:       u.Email,
		Phone:       u.Phone,
		Username:    u.Username,
		AccountName: u.Account,
		Birthday:    birthday,
		AboutMe:     u.AboutMe,
	}), nil
}

func (h *UserHandler) EditAccountName(ctx *gin.Context, req EditAccountNameReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	err := h.svc.UpdateAccountName(ctx, uc.UserId, req.AccountName)
	if errors.Is(err, user.ErrUserDuplicate) {
		return ginx.BizError("账号名已存在"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("修改成功"), nil
}

func (h *UserHandler) EditProfile(ctx *gin.Context, req EditProfileReq) (ginx.Result, error) {
	var birthday time.Time
	var err error
	if req.Birthday != "" {
		birthday, err = time.Parse(time.DateOnly, req.Birthday)
		if err != nil {
			return ginx.InvalidParamWithMsg("生日格式错误"), nil
		}
	}
	uc := jwt.MustGetUserClaims(ctx)
	err = h.svc.UpdateNonSensitiveInfo(ctx, user.User{
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

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
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
		if err = h.ClearToken(ctx); err != nil {
			h.l.WithCtx(ctx).Error("RefreshToken 一并刷新refreshToken失败", logx.Error(err))
		}
		err = h.SetLoginToken(ctx, rc.UserId)
	} else {
		// 只刷新短token
		err = h.SetJwtToken(ctx, rc.UserId, rc.Ssid)
	}
	if err != nil {
		h.l.WithCtx(ctx).Error("RefreshToken", logx.Error(err))
		ctx.JSON(http.StatusOK, ginx.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, ginx.Success())
}

func (h *UserHandler) resetPwdByCode(ctx *gin.Context, verify func(user user.User) (bool, error)) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	u, err := h.svc.Profile(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	ok, err := verify(u)
	if err != nil {
		return ginx.InternalError(), err
	}
	if !ok {
		return ginx.BizError("验证码错误"), nil
	}
	if err = h.ClearToken(ctx); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("重置密码成功, 请重新登录"), nil
}

func (h *UserHandler) ResetPwdByEmailCode(ctx *gin.Context, req EmailResetPwdReq) (ginx.Result, error) {
	return h.resetPwdByCode(ctx, func(user user.User) (bool, error) {
		if user.Email == "" {
			return false, errors.New("用户未绑定邮箱")
		}
		return h.codeSvc.Verify(ctx, resetPwdBiz, user.Email, req.Code)
	})
}

func (h *UserHandler) ChangePwd(ctx *gin.Context, req ChangePwdReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	err := h.svc.ChangePwd(ctx, uc.UserId, req.OldPassword, req.NewPassword)
	if errors.Is(err, user.ErrInvalidAccountOrPwd) {
		return ginx.BizError("原密码错误"), nil
	}
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("密码修改成功"), nil
}
