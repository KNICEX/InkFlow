package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
)

type ReviewFailRpo interface {
	Create(ctx context.Context, evt event.ReviewEvent, er error) error
	Find(ctx context.Context, offset, limit int) ([]event.ReviewEvent, error)
	Delete(ctx context.Context, ids []int64) error
}
