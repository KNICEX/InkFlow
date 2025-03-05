package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
)

type LiveInkRepo interface {
	Save(ctx context.Context, ink domain.Ink) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status domain.InkStatus) error
	FindById(ctx context.Context, id int64) (domain.Ink, error)
	ListByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]domain.Ink, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
}
