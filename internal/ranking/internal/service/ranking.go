package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/ranking/internal/domain"
)

type RankingService interface {
	RankTopN(ctx context.Context, categoryId int64, contentType int, n int) error
	TopN(ctx context.Context, categoryId int64, contentType int, offset, limit int) ([]domain.Ink, error)
}

type BatchRankingService struct {
	inkSvc ink.Service
}
