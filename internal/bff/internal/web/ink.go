package web

import (
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
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
	bizInk = "ink"

	inkPubQueue = "ink-pub-queue"
)

type InkHandler struct {
	svc            ink.Service
	inkRankService ink.RankingService
	workflowCli    client.Client
	interactiveSvc interactive.Service
	auth           middleware.Authentication
	userAggregate  *UserAggregate
	inkAggregate   *InkAggregate
	intrAggregate  *InteractiveAggregate
	l              logx.Logger
}

func NewInkHandler(svc ink.Service, userAggregate *UserAggregate, intrAggregate *InteractiveAggregate,
	intrSvc interactive.Service,
	auth middleware.Authentication,
	workflowCli client.Client, l logx.Logger) *InkHandler {
	return &InkHandler{
		svc:            svc,
		workflowCli:    workflowCli,
		interactiveSvc: intrSvc,
		auth:           auth,
		userAggregate:  userAggregate,
		inkAggregate:   NewInkAggregate(svc, userAggregate, intrAggregate),
		intrAggregate:  intrAggregate,
		l:              l,
	}
}

func (h *InkHandler) RegisterRoutes(server *gin.RouterGroup) {
	inkGroup := server.Group("/ink")

	inkGroup.GET("/detail/:id", h.auth.ExtractPayload(), ginx.Wrap(h.l, h.Detail))
	inkGroup.GET("/list", ginx.WrapBody(h.l, h.List))

	checkGroup := inkGroup.Use(h.auth.CheckLogin())
	{
		checkGroup.POST("/draft/save", ginx.WrapBody(h.l, h.SaveDraft))
		checkGroup.POST("/draft/publish/:id", ginx.Wrap(h.l, h.Publish))

		checkGroup.GET("/draft", ginx.WrapBody(h.l, h.ListDraft))
		checkGroup.GET("/pending", ginx.WrapBody(h.l, h.ListPending))
		checkGroup.GET("/private", ginx.WrapBody(h.l, h.ListPrivate))
		checkGroup.GET("/rejected", ginx.WrapBody(h.l, h.ListReviewRejected))
		checkGroup.GET("/rejected/:id", ginx.Wrap(h.l, h.DetailRejected))

		checkGroup.GET("/draft/:id", ginx.Wrap(h.l, h.DetailDraft))
		checkGroup.GET("/pending/:id", ginx.Wrap(h.l, h.DetailPending))
		checkGroup.GET("/private/:id", ginx.Wrap(h.l, h.DetailPrivate))

		checkGroup.DELETE("/draft/:id", ginx.Wrap(h.l, h.DeleteDraft))
		checkGroup.DELETE("/live/:id", ginx.Wrap(h.l, h.DeleteLive))

		checkGroup.GET("/liked", ginx.WrapBody(h.l, h.ListLiked))
		checkGroup.GET("/viewed", ginx.WrapBody(h.l, h.ListViewed))
		checkGroup.GET("/favorited", ginx.WrapBody(h.l, h.ListFavorited))
		checkGroup.POST("/withdraw/:id", ginx.Wrap(h.l, h.Withdraw))

		checkGroup.POST("/like/:id", ginx.Wrap(h.l, h.Like))
		checkGroup.DELETE("/like/:id", ginx.Wrap(h.l, h.CancelLike))
		checkGroup.POST("/favorite/:id", ginx.WrapBody(h.l, h.Favorite))
		checkGroup.DELETE("/favorite/:id", ginx.Wrap(h.l, h.CancelFavorite))

	}
}

func (h *InkHandler) SaveDraft(ctx *gin.Context, req SaveInkReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := h.svc.Save(ctx, ink.Ink{
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

func (h *InkHandler) Publish(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	err = h.svc.Publish(ctx, ink.Ink{
		Id: id,
		Author: ink.Author{
			Id: uc.UserId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}

	_, err = h.workflowCli.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        inkpub.WorkflowId(id, time.Now()),
		TaskQueue: inkPubQueue,
	}, inkpub.InkPublish, id, uc.UserId)
	if err != nil {
		h.l.WithCtx(ctx).Error("start ink publish workflow failed",
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

func (h *InkHandler) Detail(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	uc, _ := jwt.GetUserClaims(ctx)
	inkVo, err := h.inkAggregate.GetInk(ctx, id, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	go func() {
		if err = h.interactiveSvc.View(ctx.Copy(), bizInk, id, uc.UserId); err != nil {
			h.l.Error("ink handler call intrSvc view error", logx.Error(err), logx.Int64("ink_id", id))
		}
	}()
	return ginx.SuccessWithData(inkVo), nil
}

func (h *InkHandler) DetailDraft(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := h.svc.FindDraftInk(ctx, id, u.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToVO(draft)), nil
}

func (h *InkHandler) DetailPrivate(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := h.svc.FindPrivateInk(ctx, id, u.UserId)
	if err != nil {
		if errors.Is(err, ink.ErrNoPermission) {
			return ginx.NoPermission(), err
		}
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToVO(draft)), nil
}

func (h *InkHandler) DetailPending(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	uc := jwt.MustGetUserClaims(ctx)
	inkInfo, err := h.svc.FindPendingInk(ctx, id, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToVO(inkInfo)), nil
}

func (h *InkHandler) DetailRejected(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	u := jwt.MustGetUserClaims(ctx)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	draft, err := h.svc.FindRejectedInk(ctx, id, u.UserId)
	if err != nil {
		if errors.Is(err, ink.ErrNoPermission) {
			return ginx.NoPermission(), err
		}
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(inkToVO(draft)), nil
}

func (h *InkHandler) Withdraw(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.svc.Withdraw(ctx, ink.Ink{
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
func (h *InkHandler) List(ctx *gin.Context, req ListReq) (ginx.Result, error) {
	inks, err := h.svc.ListLiveByAuthorId(ctx, req.AuthorId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	if len(inks) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}

	inkIds := lo.Map(inks, func(item ink.Ink, index int) int64 {
		return item.Id
	})

	uc, _ := jwt.GetUserClaims(ctx)

	eg := errgroup.Group{}
	var author UserVO
	var intrs map[int64]InteractiveVO
	eg.Go(func() error {
		var er error
		author, er = h.userAggregate.GetUser(ctx, req.AuthorId, uc.UserId)
		return er
	})

	eg.Go(func() error {
		var er error
		intrs, er = h.intrAggregate.GetInteractiveList(ctx, bizInk, inkIds, uc.UserId)
		return er
	})

	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		inkVO := inkToVO(item)
		inkVO.Author = author
		inkVO.Interactive = intrs[item.Id]
		res = append(res, inkVO)
	}

	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) DeleteDraft(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.svc.DeleteDraft(ctx, ink.Ink{
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

func (h *InkHandler) DeleteLive(ctx *gin.Context) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.svc.DeleteLive(ctx, ink.Ink{
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

func (h *InkHandler) ListPending(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := h.svc.ListPendingByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) ListReviewRejected(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := h.svc.ListReviewRejectedByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) ListDraft(ctx *gin.Context, req ListDraftReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := h.svc.ListDraftByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) ListLiked(ctx *gin.Context, req ListMaxIdReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	// TODO 这里是不是应该把交互数据一起查出来
	likes, err := h.interactiveSvc.ListLike(ctx, bizInk, u.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(likes) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}

	inkIds := lo.Map(likes, func(item interactive.LikeRecord, index int) int64 {
		return item.BizId
	})

	inkVoMap, err := h.inkAggregate.GetInkList(ctx, inkIds, u.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(likes))
	for _, item := range likes {
		vo, ok := inkVoMap[item.BizId]
		if !ok {
			continue
		}
		res = append(res, vo)
	}

	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) ListViewed(ctx *gin.Context, req ListMaxIdReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	views, err := h.interactiveSvc.ListView(ctx, bizInk, u.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(views) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}
	inkIds := lo.Map(views, func(item interactive.ViewRecord, index int) int64 {
		return item.BizId
	})

	inkVoMap, err := h.inkAggregate.GetInkList(ctx, inkIds, u.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	inkVos := make([]InkVO, 0, len(views))
	for _, item := range views {
		vo, ok := inkVoMap[item.BizId]
		if !ok {
			continue
		}
		inkVos = append(inkVos, vo)
	}
	return ginx.SuccessWithData(inkVos), nil
}

func (h *InkHandler) ListFavorited(ctx *gin.Context, req ListFavoriteReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	fs, err := h.interactiveSvc.ListFavoriteByFid(ctx, bizInk, uc.UserId, req.Fid, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(fs) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}
	inkVOs, err := h.inkAggregate.GetInkList(ctx, lo.Map(fs, func(item interactive.FavoriteRecord, index int) int64 {
		return item.BizId
	}), uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(fs))
	for _, f := range fs {
		if vo, ok := inkVOs[f.BizId]; ok {
			res = append(res, vo)
		}
	}
	return ginx.SuccessWithData(res), nil
}
func (h *InkHandler) ListPrivate(ctx *gin.Context, req ListSelfReq) (ginx.Result, error) {
	u := jwt.MustGetUserClaims(ctx)
	inks, err := h.svc.ListPrivateByAuthorId(ctx, u.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inks))
	for _, item := range inks {
		res = append(res, inkToVO(item))
	}
	return ginx.SuccessWithData(res), nil
}

func (h *InkHandler) Like(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.interactiveSvc.Like(ctx, bizInk, id, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (h *InkHandler) CancelLike(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.interactiveSvc.CancelLike(ctx, bizInk, id, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (h *InkHandler) Favorite(ctx *gin.Context, req FavoriteReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.interactiveSvc.Favorite(ctx, bizInk, id, uc.UserId, req.FavoriteId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (h *InkHandler) CancelFavorite(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	err = h.interactiveSvc.CancelFavorite(ctx, bizInk, id, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
