package web

import (
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type InteractiveHandler struct {
	svc  interactive.Service
	auth middleware.Authentication
	l    logx.Logger
}

func NewInteractiveHandler(svc interactive.Service, auth middleware.Authentication, l logx.Logger) *InteractiveHandler {
	return &InteractiveHandler{
		svc:  svc,
		auth: auth,
		l:    l,
	}
}
func (h *InteractiveHandler) RegisterRoutes(server *gin.RouterGroup) {
	interactiveGroup := server.Group("/interactive", h.auth.CheckLogin())
	{
		interactiveGroup.GET("/favorite/:biz", ginx.Wrap(h.l, h.Favorites))
		interactiveGroup.POST("/favorite", ginx.WrapBody(h.l, h.CreateFavorite))
	}
}

func (h *InteractiveHandler) Favorites(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	biz := ctx.Param("biz")
	fs, err := h.svc.FavoriteList(ctx, biz, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(lo.Map(fs, func(item interactive.Favorite, index int) FavoriteVO {
		return favoriteToVO(item)
	})), nil
}

func (h *InteractiveHandler) CreateFavorite(ctx *gin.Context, req CreateFavReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	fid, err := h.svc.CreateFavorite(ctx, interactive.Favorite{
		Name:    req.Name,
		UserId:  uc.UserId,
		Private: req.Private,
		Biz:     req.Biz,
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(fid), nil
}
