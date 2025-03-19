package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
)

type SearchService interface {
	SearchUser(ctx context.Context, keywords []string, order domain.InkOrder) ([]domain.User, error)
	SearchUserRaw(ctx context.Context, rawExpression string) ([]domain.User, error)
	SearchInk(ctx context.Context, keywords []string, order domain.UserOrder) ([]domain.Ink, error)
	SearchInkRaw(ctx context.Context, rawExpression string) ([]domain.Ink, error)
}
