package dao

import (
	"context"
	"time"
)

type HotComment struct {
	Id        int64
	CommentId int64
	Biz       string `gorm:"index:biz_type"`
	BizId     int64  `gorm:"index:biz_type"`
	Hot       int64  `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HotCommentDAO interface {
	Insert(ctx context.Context, c []HotComment) error
	FindByBiz(ctx context.Context, biz string, bizId int64, offset int, limit int) ([]HotComment, error)
}
