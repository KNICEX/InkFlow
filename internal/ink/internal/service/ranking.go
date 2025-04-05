package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
)

type RankingService interface {
	TopN(ctx context.Context, n int) error
	FindTopN(ctx context.Context, offset int, limit int) ([]domain.Ink, error)
}
