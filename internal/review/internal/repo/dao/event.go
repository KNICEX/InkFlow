package dao

import (
	"context"
	"database/sql"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

type ReviewFailDAO interface {
	Insert(ctx context.Context, event ReviewFail) error
	Find(ctx context.Context, offset, limit int) ([]ReviewFail, error)
	Delete(ctx context.Context, ids []int64) error
}

type ReviewFail struct {
	Id         int64
	WorkflowId string
	Event      string
	Error      string
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
	// DeletedAt 软删除做为错误分析
	DeletedAt sql.NullTime
}

type GormReviewFailDAO struct {
	db   *gorm.DB
	node snowflakex.Node
}

func NewGormReviewFailDAO(db *gorm.DB, node snowflakex.Node) ReviewFailDAO {
	return &GormReviewFailDAO{
		db:   db,
		node: node,
	}
}

func (g *GormReviewFailDAO) Insert(ctx context.Context, event ReviewFail) error {
	event.Id = g.node.NextID()
	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt
	return g.db.WithContext(ctx).Create(&event).Error
}

func (g *GormReviewFailDAO) Find(ctx context.Context, offset, limit int) ([]ReviewFail, error) {
	var events []ReviewFail
	err := g.db.WithContext(ctx).Where("deleted_at IS NULL").
		Offset(offset).Limit(limit).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (g *GormReviewFailDAO) Delete(ctx context.Context, ids []int64) error {
	return g.db.WithContext(ctx).Where("id IN ?", ids).Update("deleted_at", time.Now()).Error
}
