package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/comment/internal/domain"
)

type CommentService interface {
	// LoadLastedCommentList 倒序加载(最新评论)
	LoadLastedCommentList(ctx context.Context, biz string, bizId int64, maxId int64, limit int)
	LoadHotCommentList(ctx context.Context, biz string, bizId int64, offset int, limit int)
	DeleteComment(ctx context.Context, id int64) error
	CreateComment(ctx context.Context, comment domain.Comment) error
	// LoadMoreRepliesByRid 根据rootId加载所有子评论
	LoadMoreRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int) ([]domain.Comment, error)
	// LoadMoreRepliesByPid 根据parentId加载所有子评论
	LoadMoreRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int) ([]domain.Comment, error)
}
