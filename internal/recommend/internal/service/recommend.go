package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
)

type RecommendService interface {
	FindSimilarInk(ctx context.Context, inkId int64) ([]domain.Ink, error)
	FindSimilarUser(ctx context.Context, userId int64) ([]domain.User, error)
	FindSimilarAuthor(ctx context.Context, authorId int64) ([]domain.User, error)
	FindPopular(ctx context.Context, offset, limit int) ([]domain.Ink, error)
	FindRecommendInk(ctx context.Context, userId int64, offset, limit int) ([]domain.Ink, error)
	FindRecommendAuthor(ctx context.Context, userId int64, offset, limit int) ([]domain.User, error)
}
