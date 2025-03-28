package service

import (
	"context"
)

type RecommendService interface {
	FindSimilarInk(ctx context.Context, inkId int64, offset, limit int) ([]int64, error)
	FindSimilarUser(ctx context.Context, userId int64, offset, limit int) ([]int64, error)
	FindSimilarAuthor(ctx context.Context, authorId int64, offset, limit int) ([]int64, error)
	FindPopular(ctx context.Context, offset, limit int) ([]int64, error)
	FindRecommendInk(ctx context.Context, userId int64, offset, limit int) ([]int64, error)
	FindRecommendAuthor(ctx context.Context, userId int64, offset, limit int) ([]int64, error)
}
