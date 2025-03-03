package dao

import "context"

type DraftDAO interface {
	Insert(ctx context.Context, d DraftInk) error
	GetById(ctx context.Context, id int64) (DraftInk, error)
	UpdateByIdAnAuthorId(ctx context.Context, d DraftInk) error
	FindByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]DraftInk, error)
}
