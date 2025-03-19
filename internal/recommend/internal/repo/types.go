package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
)

type RecommendRepo interface {
	AddInk(ctx context.Context, ink domain.Ink)
	HideInk(ctx context.Context, inkId int64)
	UnHideInk(ctx context.Context, inkId int64)
	AddFeedback(ctx context.Context, feedback domain.Feedback)
	AddUser(ctx context.Context, user domain.User)
}
