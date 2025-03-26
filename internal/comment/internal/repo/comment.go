package repo

import (
	"context"

	"github.com/KNICEX/InkFlow/internal/comment/internal/domain"
)

// CommentRepo defines the data access operations for comments
type CommentRepo interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DelComment(ctx context.Context, commentId int64, uid int64) error
	FindByBiz(ctx context.Context, biz string, bizId int64, uid, maxId int64, limit int) ([]domain.Comment, error)
	FindByRootId(ctx context.Context, rootId int64, uid, maxId int64, limit int) ([]domain.Comment, error)
	FindByParentId(ctx context.Context, parentId int64, uid, maxId int64, limit int) ([]domain.Comment, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Comment, error)
}
