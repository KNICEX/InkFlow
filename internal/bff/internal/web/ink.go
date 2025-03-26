package web

import (
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/internal/workflow/inkpub"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.temporal.io/sdk/client"
	"golang.org/x/sync/errgroup"
	"strconv"
	"time"
)

const (
	inkBiz = "ink"

	inkPubQueue = "ink-pub-queue"
)

type InkHandler struct {
	svc            ink.Service
	userSvc        user.Service
	workflowCli    client.Client
	interactiveSvc interactive.Service
	followService  relation.FollowService
	auth           middleware.Authentication
	l              logx.Logger
}

func NewInkHandler(svc ink.Service, userSvc user.Service, interactiveSvc interactive.Service,
	followService relation.FollowService, auth middleware.Authentication,
	workflowCli client.Client, l logx.Logger) *InkHandler {
	return &InkHandler{
		svc:            svc,
		userSvc:        userSvc,
		followService:  followService,
		workflowCli:    workflowCli,
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
		checkGroup.POST("/draft/publish/:id", ginx.Wrap(handler.l, handler.Publish))
		checkGroup.GET("/draft/detail/:id", ginx.Wrap(handler.l, handler.DetailDraft))
		checkGroup.POST("/draft/delete/:id", ginx.Wrap(handler.l, handler.DeleteDraft))
		checkGroup.GET("/private/detail/:id", ginx.Wrap(handler.l, handler.DetailPrivate))
		checkGroup.GET("/review/detail/:id", ginx.Wrap(handler.l, handler.Detail))
		checkGroup.POST("/live/delete/:id", ginx.Wrap(handler.l, handler.DeleteLive))
		checkGroup.POST("/draft/list", ginx.WrapBody(handler.l, handler.ListDraft))
		checkGroup.POST("/pending/list", ginx.WrapBody(handler.l, handler.ListPending))
		checkGroup.POST("/private/list", ginx.WrapBody(handler.l, handler.ListPrivate))
		checkGroup.POST("/rejected/list", ginx.WrapBody(handler.l, handler.ListReviewRejected))
		checkGroup.POST("/liked/list", ginx.WrapBody(handler.l, handler.ListLiked))
		checkGroup.POST("/viewed/list", ginx.WrapBody(handler.l, handler.ListViewed))
		checkGroup.POST("/withdraw/:id", ginx.Wrap(handler.l, handler.Withdraw))
	}
}

func (handler *InkHandler) SaveDraft(ctx *gin.Context, req SaveInkReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := handler.svc.Save(ctx, ink.Ink{
		Id:          req.Id,
		Title:       req.Title,
		Cover:       req.Cover,
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

func (handler *InkHandler) Publish(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	id, err = handler.svc.Publish(ctx, ink.Ink{
		Id: id,
		Author: ink.Author{
			Id: u.UserId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}

	_, err = handler.workflowCli.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        inkpub.WorkflowId(id, time.Now()),
		TaskQueue: inkPubQueue,
	}, inkpub.InkPublish, id)
	if err != nil {
		handler.l.WithCtx(ctx).Error("start ink publish workflow failed",
			logx.Int64("inkId", id),
			logx.Error(err))
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
	inkDetail, err := handler.svc.FindLiveInk(ctx, id)
	if err != nil {
		return ginx.InternalError(), nil
	}

	// 无所谓是否登录，id为0就是没登录
	readUser, _ := jwt.GetUserClaims(ctx)
	readUserId := readUser.UserId

	eg := errgroup.Group{}
	var author user.User
	var followInfo relation.FollowStatistic
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
		var er error
		followInfo, er = handler.followService.FindFollowStats(ctx, inkDetail.Author.Id, readUserId)
		return er
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

	authorProfile := userToUserVO(author)
	authorProfile.Followers = followInfo.Followers
	authorProfile.Following = followInfo.Following
	authorProfile.Followed = followInfo.Followed
	return ginx.SuccessWithData(InkVO{
		InkBaseVO:   inkToInkBaseVO(inkDetail),
		Author:      authorProfile,
		Interactive: intrToVo(intr),
	}), nil
}

func (handler *InkHandler) DetailDraft(ctx *gin.Context) (ginx.Result, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := handler.svc.FindDraftInk(ctx, id, u.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToInkBaseVO(draft)), nil
}

func (handler *InkHandler) DetailPrivate(ctx *gin.Context) (ginx.Result, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := handler.svc.FindPrivateInk(ctx, id, u.UserId)
	if err != nil {
		if errors.Is(err, ink.ErrNoPermission) {
			return ginx.NoPermission(), err
		}
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToInkBaseVO(draft)), nil
}

func (handler *InkHandler) DetailRejected(ctx *gin.Context) (ginx.Result, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := handler.svc.FindRejectedInk(ctx, id, u.UserId)
	if err != nil {
		if errors.Is(err, ink.ErrNoPermission) {
			return ginx.NoPermission(), err
		}
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToInkBaseVO(draft)), nil
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

	if len(inks) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
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

	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		author, ok := users[item.Author.Id]
		if !ok {
			continue
		}
		intr, ok := intrs[item.Id]
		if !ok {
			continue
		}
		res = append(res, InkVO{
			InkBaseVO:   inkToInkBaseVO(item),
			Author:      userToUserVO(author),
			Interactive: intrToVo(intr),
		})
	}

	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) DeleteDraft(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = handler.svc.DeleteDraft(ctx, ink.Ink{
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

func (handler *InkHandler) DeleteLive(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = handler.svc.DeleteLive(ctx, ink.Ink{
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

func (handler *InkHandler) ListPending(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := handler.svc.ListPendingByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkBaseVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToInkBaseVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) ListReviewRejected(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := handler.svc.ListReviewRejectedByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkBaseVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToInkBaseVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) ListDraft(ctx *gin.Context, req ListDraftReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := handler.svc.ListDraftByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkBaseVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToInkBaseVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) ListLiked(ctx *gin.Context, req ListMaxIdReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	// TODO 这里是不是应该把交互数据一起查出来
	likes, err := handler.interactiveSvc.ListLike(ctx, inkBiz, u.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	inkIds := lo.Map(likes, func(item interactive.LikeRecord, index int) int64 {
		return item.BizId
	})

	eg := errgroup.Group{}
	var inkMap map[int64]ink.Ink
	var authorMap map[int64]user.User
	eg.Go(func() error {
		var er error
		inkMap, er = handler.svc.FindByIds(ctx, inkIds)
		return er
	})
	eg.Go(func() error {
		var er error
		authorMap, er = handler.userSvc.FindByIds(ctx, inkIds)
		return er
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	res := make([]InkVO, 0, len(likes))
	for _, item := range likes {
		i, ok := inkMap[item.BizId]
		if !ok {
			continue
		}
		author, ok := authorMap[item.BizId]
		if !ok {
			continue
		}
		res = append(res, InkVO{
			InkBaseVO: inkToInkBaseVO(i),
			Author:    userToUserVO(author),
		})
	}

	return ginx.SuccessWithData(res), nil
}

func (handler *InkHandler) ListViewed(ctx *gin.Context, req ListMaxIdReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	views, err := handler.interactiveSvc.ListView(ctx, inkBiz, u.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	inkIds := lo.Map(views, func(item interactive.ViewRecord, index int) int64 {
		return item.BizId
	})

	eg := errgroup.Group{}
	var inkMap map[int64]ink.Ink
	var authorMap map[int64]user.User
	eg.Go(func() error {
		var er error
		inkMap, er = handler.svc.FindByIds(ctx, inkIds)
		return er
	})
	eg.Go(func() error {
		var er error
		authorMap, er = handler.userSvc.FindByIds(ctx, inkIds)
		return er
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	res := make([]InkVO, 0, len(views))
	for _, item := range views {
		i, ok := inkMap[item.BizId]
		if !ok {
			continue
		}
		author, ok := authorMap[item.BizId]
		if !ok {
			continue
		}
		res = append(res, InkVO{
			InkBaseVO: inkToInkBaseVO(i),
			Author:    userToUserVO(author),
		})
	}
	return ginx.SuccessWithData(res), nil
}
func (handler *InkHandler) ListPrivate(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := handler.svc.ListPrivateByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkBaseVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToInkBaseVO(item))
	}
	return ginx.SuccessWithData(res), nil
}
