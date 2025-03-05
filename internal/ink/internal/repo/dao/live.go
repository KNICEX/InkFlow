package dao

import "context"

type LiveDAO interface {
	Upsert(ctx context.Context, d LiveInk) error
	UpdateStatus(ctx context.Context, inkId int64, authorId int64, status int) error
	FindById(ctx context.Context, id int64) (LiveInk, error)
	FindByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]LiveInk, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]LiveInk, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]LiveInk, error)
}
