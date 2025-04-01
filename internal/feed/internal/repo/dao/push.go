package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

type PushFeed struct {
	Id        int64
	UserId    int64  `gorm:"index:idx_uid_biz_status"`
	Biz       string `gorm:"index:idx_uid_biz_status"`
	BizId     int64
	Content   string
	Status    int `gorm:"index:idx_uid_biz_status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PushFeedDAO interface {
	CreatePush(ctx context.Context, push []PushFeed) error
	FindHidden(ctx context.Context, limit int) ([]PushFeed, error)
	BatchDelete(ctx context.Context, ids []int64) error
	UpdateStatus(ctx context.Context, feed PushFeed) error
	FindPush(ctx context.Context, uid int64, maxId, timestamp int64, limit int) ([]PushFeed, error)
	FindPushByBiz(ctx context.Context, uid int64, biz string, maxId, timestamp int64, limit int) ([]PushFeed, error)
}

type GormPushFeedDAO struct {
	db   *gorm.DB
	node snowflakex.Node
}

func NewGormPushFeedDAO(db *gorm.DB, node snowflakex.Node) PushFeedDAO {
	return &GormPushFeedDAO{
		db:   db,
		node: node,
	}
}

func (dao *GormPushFeedDAO) CreatePush(ctx context.Context, push []PushFeed) error {
	for i := range push {
		push[i].Id = dao.node.NextID()
	}
	return dao.db.WithContext(ctx).Create(&push).Error
}

func (dao *GormPushFeedDAO) FindHidden(ctx context.Context, limit int) ([]PushFeed, error) {
	tx := dao.db.WithContext(ctx).Model(&PushFeed{}).
		Where("status = ?", FeedStatusHidden).
		Limit(limit)
	var pushes []PushFeed
	err := tx.Find(&pushes).Error
	if err != nil {
		return nil, err
	}
	return pushes, nil
}

func (dao *GormPushFeedDAO) BatchDelete(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Model(&PushFeed{}).
		Where("id IN ?", ids).
		Delete(&PushFeed{}).Error
}

func (dao *GormPushFeedDAO) UpdateStatus(ctx context.Context, feed PushFeed) error {
	feed.UpdatedAt = time.Now()
	return dao.db.WithContext(ctx).Model(&PushFeed{}).
		Where("biz = ? AND feed_id = ?", feed.Biz, feed.BizId).
		Update("status", feed.Status).Error
}

func (dao *GormPushFeedDAO) FindPush(ctx context.Context, uid int64, maxId, timestamp int64, limit int) ([]PushFeed, error) {
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		tx = tx.Where("user_id = ? AND status = ?", uid, FeedStatusNormal)
	} else {
		tx = tx.Where("user_id = ? AND id < ? AND created_at < ? AND status = ?", uid, maxId, time.UnixMilli(timestamp), FeedStatusNormal)
	}
	var feeds []PushFeed
	err := tx.Order("id desc").Limit(limit).Find(&feeds).Error
	return feeds, err
}

func (dao *GormPushFeedDAO) FindPushByBiz(ctx context.Context, uid int64, biz string, maxId, timestamp int64, limit int) ([]PushFeed, error) {
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		tx = tx.Where("user_id = ? AND biz = ? AND status = ?", uid, biz, FeedStatusNormal)
	} else {
		tx = tx.Where("user_id = ? AND biz = ? AND id < ? AND created_at < ? AND status = ?", uid, biz, maxId, time.UnixMilli(timestamp), FeedStatusNormal)
	}
	var feeds []PushFeed
	err := tx.Order("id desc").Limit(limit).Find(&feeds).Error
	return feeds, err

}
