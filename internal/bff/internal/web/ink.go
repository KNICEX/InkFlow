package web

import (
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"strconv"
)

const (
	inkBiz = "ink"
)

type InkHandler struct {
	svc            ink.Service
	userSvc        user.Service
	interactiveSvc interactive.Service
	auth           middleware.Authentication
	l              logx.Logger
}

func NewInkHandler(svc ink.Service, userSvc user.Service, interactiveSvc interactive.Service, auth middleware.Authentication, l logx.Logger) *InkHandler {
	return &InkHandler{
		svc:            svc,
		userSvc:        userSvc,
		interactiveSvc: interactiveSvc,
		auth:           auth,
		l:              l,
	}
}

func (handler *InkHandler) RegisterRoutes(server *gin.RouterGroup) {
	inkGroup := server.Group("/ink")

	inkGroup.GET("/detail/:id", ginx.Wrap(handler.l, handler.Detail))
	inkGroup.POST("/list", ginx.WrapBody(handler.l, handler.List))

	checkGroup := inkGroup.Use(handler.auth.CheckLogin())
	{
		checkGroup.POST("/draft/save", ginx.WrapBody(handler.l, handler.SaveDraft))
		checkGroup.POST("/draft/publish", ginx.WrapBody(handler.l, handler.Publish))
		checkGroup.GET("/draft/detail/:id", ginx.Wrap(handler.l, handler.DraftDetail))
		checkGroup.POST("/draft/list", ginx.WrapBody(handler.l, handler.List))
		checkGroup.POST("/withdraw/:id", ginx.Wrap(handler.l, handler.Withdraw))
	}
}

func (handler *InkHandler) SaveDraft(ctx *gin.Context, req SaveInkReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := handler.svc.Save(ctx, ink.Ink{
		Id:          req.Id,
		Title:       req.Title,
		Summary:     req.Summary,
		ContentHtml: req.ContentHtml,
		ContentMeta: req.ContentMeta,
		Tags:        req.Tags,
		Author: ink.Author{
			Id: u.UserId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	type SaveResp struct {
		Id int64 `json:"id"`
	}
	return ginx.SuccessWithData(SaveResp{
		Id: id,
	}), nil
}

func (handler *InkHandler) Publish(ctx *gin.Context, req PublishInkReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := handler.svc.Publish(ctx, ink.Ink{
		Id: req.Id,
		Author: ink.Author{
			Id: u.UserId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	type PublishResp struct {
		Id int64 `json:"id"`
	}
	return ginx.SuccessWithData(PublishResp{
		Id: id,
	}), nil
}

func (handler *InkHandler) Detail(ctx *gin.Context) (ginx.Result, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	inkDetail, err := handler.svc.GetLiveInk(ctx, id)
	if err != nil {
		return ginx.InternalError(), nil
	}

	// 无所谓是否登录，id为0就是没登录
	readUser, _ := jwt.GetUserClaims(ctx)
	readUserId := readUser.UserId

	eg := errgroup.Group{}
	var author user.User
	var intr interactive.Interactive
	eg.Go(func() error {
		var er error
		author, er = handler.userSvc.FindById(ctx, inkDetail.Author.Id)
		return er
	})
	eg.Go(func() error {
		var er error
		intr, er = handler.interactiveSvc.Get(ctx, inkBiz, inkDetail.Id, readUserId)
		return er
	})

	eg.Go(func() error {
		// TODO 获取用户关注信息
		return nil
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	go func() {
		er := handler.interactiveSvc.View(ctx, inkBiz, inkDetail.Id, readUserId)
		if er != nil {
			handler.l.WithCtx(ctx).Error("send read event failed when get live article",
				logx.Int64("user_id", readUserId),
				logx.Int64("ink_id", inkDetail.Id),
				logx.Error(er))
		}
	}()

	authorProfile := UserProfileFromDomain(author)
	return ginx.SuccessWithData(InkDetailResp{
		InkBaseInfo: InkBaseInfoFromDomain(inkDetail),
		Author:      authorProfile,
		Interactive: InteractiveVOFromDomain(intr),
	}), nil
}

func (handler *InkHandler) DraftDetail(ctx *gin.Context) (ginx.Result, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := handler.svc.GetDraftInk(ctx, u.UserId, id)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(InkBaseInfoFromDomain(draft)), nil
}

func (handler *InkHandler) Withdraw(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = handler.svc.Withdraw(ctx, ink.Ink{
		Id: id,
		Author: ink.Author{
			Id: u.UserId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
func (handler *InkHandler) List(ctx *gin.Context, req ListReq) (ginx.Result, error) {
	inks, err := handler.svc.ListLiveByAuthorId(ctx, req.AuthorId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uids := lo.Map(inks, func(item ink.Ink, index int) int64 {
		return item.Author.Id
	})
	inkIds := lo.Map(inks, func(item ink.Ink, index int) int64 {
		return item.Id
	})

	readUser, _ := jwt.GetUserClaims(ctx)
	readUserId := readUser.UserId

	eg := errgroup.Group{}
	var users map[int64]user.User
	var intrs map[int64]interactive.Interactive
	eg.Go(func() error {
		var er error
		users, er = handler.userSvc.FindByIds(ctx, uids)
		return er
	})

	eg.Go(func() error {
		var er error
		intrs, er = handler.interactiveSvc.GetMulti(ctx, inkBiz, inkIds, readUserId)
		return er
	})

	eg.Go(func() error {
		// TODO 批量获取用户关注信息
		return nil
	})

	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	res := make([]InkDetailResp, 0, len(inks))
	for _, item := range inks {
		author, ok := users[item.Author.Id]
		if !ok {
			continue
		}
		intr, ok := intrs[item.Id]
		if !ok {
			continue
		}
		res = append(res, InkDetailResp{
			InkBaseInfo: InkBaseInfoFromDomain(item),
			Author:      UserProfileFromDomain(author),
			Interactive: InteractiveVOFromDomain(intr),
		})
	}

	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) ListDraft(ctx *gin.Context, req ListDraftReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := handler.svc.ListDraftByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkBaseInfo, 0, len(inks))
	for _, item := range inks {
		res = append(res, InkBaseInfoFromDomain(item))
	}
	return ginx.SuccessWithData(res), nil
}
