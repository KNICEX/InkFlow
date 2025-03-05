package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
)

type DraftInkRepo interface {
	Create(ctx context.Context, ink domain.Ink) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Ink, error)
	Update(ctx context.Context, ink domain.Ink) error
	ListByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]domain.Ink, error)
}
