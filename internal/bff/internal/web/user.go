package web

import (
	"errors"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/jwtx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	loginBiz    = "login"
	resetPwdBiz = "reset_pwd"
)

const (
	registerJwtKey = "register"
)

type UserHandler struct {
	svc           user.Service
	codeSvc       code.Service
	followService relation.FollowService
	phoneReg      *regexp.Regexp
	emailReg      *regexp.Regexp
	l             logx.Logger
	auth          middleware.Authentication
	jwt.Handler
}

func NewUserHandler(svc user.Service,
	codeSvc code.Service, followService relation.FollowService,
	jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) *UserHandler {
	return &UserHandler{
		svc:           svc,
		codeSvc:       codeSvc,
		followService: followService,
		phoneReg:      regexp.MustCompile(`^1[3456789]\d{9}$`),
		emailReg:      regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`),
		Handler:       jwtHandler,
		l:             log,
		auth:          auth,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.RouterGroup) {
	userGroup := server.Group("/user")
	// 登录验证码
	userGroup.POST("/verify/send/login", ginx.WrapBody(h.l, h.SendLoginCode))

	loginGroup := userGroup.Group("/login")
	{
		// 邮箱验证码登录(未注册会返回临时凭证)
		loginGroup.POST("/email", ginx.WrapBody(h.l, h.LoginEmail))
		// 邮箱密码登录
		loginGroup.POST("/pwd/email", ginx.WrapBody(h.l, h.LoginEmailPwd))
		// 账号密码登录
		loginGroup.POST("/pwd/account", ginx.WrapBody(h.l, h.LoginAccountPwd))
	}
	registerGroup := userGroup.Group("/register")
	{
		registerGroup.POST("/email", ginx.WrapBody(h.l, h.RegisterByEmail))
	}

	// 刷新token
	userGroup.POST("/refresh_token", h.RefreshToken)
	// 获取用户基础信息
	userGroup.POST("/profile", h.auth.ExtractPayload(), ginx.WrapBody(h.l, h.Profile))

	// 需要登录
	checkGroup := userGroup.Group("")
	checkGroup.Use(h.auth.CheckLogin())
	{
		checkGroup.GET("/logout", ginx.Wrap(h.l, h.Logout))

		// 修改个人信息
		checkGroup.PUT("/profile", ginx.WrapBody(h.l, h.EditProfile))

		// 发送重置密码验证码
		//checkGroup.POST("/verify/send/reset/sms", ginx.Wrap(h.l, h.SendResetPwdCodeSms))
		checkGroup.POST("/verify/send/reset/email", ginx.Wrap(h.l, h.SendResetPwdCodeEmail))

		// 重置密码
		checkGroup.POST("/pwd/reset/email", ginx.WrapBody(h.l, h.ResetPwdByEmailCode))
		//checkGroup.POST("/pwd/reset/sms", ginx.WrapBody(h.l, h.ResetPwdBySmsCode))
		checkGroup.POST("/pwd/reset/old", ginx.WrapBody(h.l, h.ChangePwd))

		{
			// 关注
			checkGroup.POST("/follow/:id", ginx.Wrap(h.l, h.Follow))
			// 取消关注
			checkGroup.DELETE("/follow/:id", ginx.Wrap(h.l, h.CancelFollow))
			// 关注列表
			checkGroup.GET("/follow/:id/following", ginx.WrapBody(h.l, h.FollowingList))
			// 粉丝列表
			checkGroup.GET("/follow/:id/follower", ginx.WrapBody(h.l, h.FollowerList))
		}
	}

}

func (h *UserHandler) sendCodeWithSvc(ctx *gin.Context, biz, recipient string) (ginx.Result, error) {
	err := h.codeSvc.Send(ctx, biz, recipient)
	switch {
	case err == nil:
		return ginx.SuccessWithMsg("验证码发送成功"), nil
	case errors.Is(err, code.ErrCodeSendTooMany):
		// 这里warn一下就可以了
		h.l.WithCtx(ctx).Warn("验证码发送太频繁", logx.String("recipient", recipient))
		return ginx.BizError("验证码发送太频繁"), nil
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
	if err != nil {
		if errors.Is(err, code.ErrCodeVerifyLimit) {
			h.l.WithCtx(ctx).Warn("验证码验证太频繁", logx.String("recipient", recipient))
			return ginx.BizError("验证码验证太频繁"), nil
		} else {
			return ginx.InternalError(), err
		}
	}
	if !ok {
		return ginx.BizError("验证码错误"), nil
	}
	return ginx.Success(), nil
}

func (h *UserHandler) SendResetPwdCodeEmail(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	u, err := h.svc.FindById(ctx, uc.UserId)
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

	u, err = h.svc.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			// 进入两步注册,
			// TODO 潜在的并发可能，同时存在两个凭证，可能需要借助redis做唯一验证
			token, er := jwtx.Generate(req.Email, time.Minute*10, registerJwtKey)
			if er != nil {
				h.l.WithCtx(ctx).Error("生成注册token失败", logx.Error(er))
				return ginx.InternalError(), er
			}
			return ginx.SuccessWithData(token), nil
		}
		return ginx.InternalError(), err
	}

	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("登录成功"), nil
}

func (h *UserHandler) RegisterByEmail(ctx *gin.Context, req RegisterByEmailReq) (ginx.Result, error) {
	email, err := jwtx.Parse[string](req.Token, registerJwtKey)
	if err != nil {
		return ginx.BizError("无效的凭证"), nil
	}
	if email != req.Email {
		return ginx.BizError("邮箱不匹配"), nil
	}
	u, err := h.svc.Create(ctx, user.User{
		Email:    req.Email,
		Username: req.Username,
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrUserDuplicate) {
			// TODO 当然也可能是邮箱重复 (概率很低，前端会预检账号名)
			return ginx.BizError("账户名重复"), nil
		}
		return ginx.InternalError(), err
	}
	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithMsg("注册成功"), nil
}

// LoginEmailPwd 邮箱密码登录
func (h *UserHandler) LoginEmailPwd(ctx *gin.Context, req LoginEmailPwdReq) (ginx.Result, error) {
	u, err := h.svc.LoginEmailPwd(ctx, req.Email, req.Password)
	if errors.Is(err, user.ErrInvalidAccountOrPwd) || errors.Is(err, user.ErrUserNotFound) {
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

func (h *UserHandler) Profile(ctx *gin.Context, req ProfileReq) (ginx.Result, error) {
	uc, logined := jwt.GetUserClaims(ctx)
	var u user.User
	var follow relation.FollowStatistic
	eg := errgroup.Group{}
	var fromToken bool

	if req.Uid > 0 {
		// 根据id查询
		eg.Go(func() error {
			var er error
			u, er = h.svc.FindById(ctx, req.Uid)
			return er
		})
		eg.Go(func() error {
			var er error
			follow, er = h.followService.FindFollowStats(ctx, req.Uid, uc.UserId)
			return er
		})
	} else if req.Account != "" {
		// 根据账号查询
		eg.Go(func() error {
			var er error
			u, er = h.svc.FindByAccount(ctx, req.Account)
			if er != nil {
				return er
			}
			follow, er = h.followService.FindFollowStats(ctx, u.Id, uc.UserId)
			if er != nil {
				return er
			}
			return er
		})
	} else if !logined {
		// 未登录查询自己的信息
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ginx.BizError("未登录"), nil
	} else {
		// 直接通过token凭证查询自己的信息
		fromToken = true
		eg.Go(func() error {
			var er error
			u, er = h.svc.FindById(ctx, uc.UserId)
			return er
		})
		eg.Go(func() error {
			var er error
			follow, er = h.followService.FindFollowStats(ctx, req.Uid, uc.UserId)
			return er
		})
	}

	err := eg.Wait()
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			if fromToken {
				// 合法凭证，但是用户不存在
				h.l.Error("user not found", logx.Error(err), logx.Int64("uid", uc.UserId))
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return ginx.BizError("凭证失效"), nil
			}
			return ginx.BizError("用户不存在"), nil
		}
		return ginx.InternalError(), err
	}

	res := userToVO(u)
	res.Followers = follow.Followers
	res.Following = follow.Following
	res.Followed = follow.Followed
	return ginx.SuccessWithData(userToVO(u)), nil
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
		Username: req.Username,
		Avatar:   req.Avatar,
		Banner:   req.Banner,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
		Links:    req.Links,
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
	rc, err := jwt.ParseRefreshClaims(req.RefreshToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ok, err := h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		h.l.WithCtx(ctx).Error("RefreshToken", logx.Error(err))
		ctx.JSON(http.StatusOK, ginx.InternalError())
		return
	}

	if !ok {
		// refresh token 过期了
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
	u, err := h.svc.FindById(ctx, uc.UserId)
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

func (h *UserHandler) Follow(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), nil
	}
	err = h.followService.Follow(ctx, uc.UserId, id)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
func (h *UserHandler) CancelFollow(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), nil
	}
	err = h.followService.CancelFollow(ctx, uc.UserId, id)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (h *UserHandler) FollowingList(ctx *gin.Context, req FollowListReq) (ginx.Result, error) {
	return h.followList(ctx, req, true)
}

func (h *UserHandler) FollowerList(ctx *gin.Context, req FollowListReq) (ginx.Result, error) {
	return h.followList(ctx, req, false)
}

func (h *UserHandler) followList(ctx *gin.Context, req FollowListReq, following bool) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), nil
	}
	var follows []relation.FollowStatistic
	if following {
		follows, err = h.followService.FollowingList(ctx, id, uc.UserId, req.MaxId, req.Limit)
	} else {
		follows, err = h.followService.FollowerList(ctx, id, uc.UserId, req.MaxId, req.Limit)
	}
	uids := lo.Map(follows, func(item relation.FollowStatistic, index int) int64 {
		return item.Uid
	})
	users, err := h.svc.FindByIds(ctx, uids)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]UserVO, 0, len(users))
	for _, v := range follows {
		if u, ok := users[v.Uid]; ok {
			profile := userToVO(u)
			profile.Followers = v.Followers
			profile.Following = v.Following
			profile.Followed = v.Followed
			res = append(res, profile)
		}
	}
	return ginx.SuccessWithData(res), nil
}
