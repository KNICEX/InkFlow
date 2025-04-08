package web

import (
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	inkRankService ink.RankingService
	l              logx.Logger
}

func NewStatsHandler(inkRankService ink.RankingService, l logx.Logger) *StatsHandler {
	return &StatsHandler{
		inkRankService: inkRankService,
		l:              l,
	}
}

func (h *StatsHandler) RegisterRoutes(server *gin.RouterGroup) {
	statsGroup := server.Group("/stats")
	{
		statsGroup.GET("/top-tags", ginx.WrapBody(h.l, h.TopTags))
	}
}

func (h *StatsHandler) TopTags(ctx *gin.Context, req OffsetPagedReq) (ginx.Result, error) {
	tags, err := h.inkRankService.FindTopNTag(ctx, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(tags), nil
}
