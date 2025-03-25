package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

type PullFeed struct {
	Id        int64
	UserId    int64  `gorm:"index:idx_uid_type_status"`
	FeedType  string `gorm:"index:idx_uid_type_status"`
	FeedId    int64
	Content   string
	Status    int `gorm:"index:idx_uid_type_status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	FeedStatusNormal = 1
	FeedStatusHidden = 2
)

type PullFeedDAO interface {
	CreatePull(ctx context.Context, pull PullFeed) error
	FindHidden(ctx context.Context, limit int) ([]PullFeed, error)
	BatchDelete(ctx context.Context, ids []int64) error
	UpdateStatus(ctx context.Context, feed PullFeed) error
	FindPull(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]PullFeed, error)
	FindPullByType(ctx context.Context, uids []int64, feedType string, maxId, timestamp int64, limit int) ([]PullFeed, error)
}

type GormFeedPullDAO struct {
	node snowflakex.Node
	db   *gorm.DB
}

func (dao *GormFeedPullDAO) CreatePull(ctx context.Context, pull PullFeed) error {
	pull.Id = dao.node.NextID()
	err := dao.db.WithContext(ctx).Create(&pull).Error
	err, _ = gormx.CheckDuplicateErr(err)
	return err
}

func (dao *GormFeedPullDAO) FindHidden(ctx context.Context, limit int) ([]PullFeed, error) {
	tx := dao.db.WithContext(ctx).Model(&PullFeed{}).
		Where("status = ?", FeedStatusHidden).
		Limit(limit)
	var pulls []PullFeed
	err := tx.Find(&pulls).Error
	if err != nil {
		return nil, err
	}
	return pulls, nil
}

func (dao *GormFeedPullDAO) BatchDelete(ctx context.Context, ids []int64) error {
	err := dao.db.WithContext(ctx).Model(&PullFeed{}).
		Where("id in ?", ids).
		Delete(&PullFeed{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *GormFeedPullDAO) UpdateStatus(ctx context.Context, feed PullFeed) error {
	err := dao.db.WithContext(ctx).Model(&PullFeed{}).
		Where("user_id = ? AND feed_type = ? AND feed_id = ?", feed.UserId, feed.FeedType, feed.FeedId).
		Update("status", feed.Status).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *GormFeedPullDAO) FindPull(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]PullFeed, error) {
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		tx = tx.Where("user_id in ? AND status = ?", uids, FeedStatusNormal)
	} else {
		tx = tx.Where("user_id in ? AND id < ? AND created_at < ? AND status = ?", uids, maxId, time.UnixMilli(timestamp), FeedStatusNormal)
	}
	feeds := make([]PullFeed, 0, limit)
	err := tx.Order("id desc").Limit(limit).Find(&feeds).Error
	return feeds, err
}

func (dao *GormFeedPullDAO) FindPullByType(ctx context.Context, uids []int64, feedType string, maxId, timestamp int64, limit int) ([]PullFeed, error) {
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		tx = tx.Where("user_id in ? AND feed_type = ? AND status = ?", uids, feedType, FeedStatusNormal)
	} else {
		tx = tx.Where("user_id in ? AND feed_type = ? AND id < ? AND created_at < ? AND status = ?", uids, feedType, maxId, time.UnixMilli(timestamp), FeedStatusNormal)
	}
	var feeds []PullFeed
	err := tx.Order("id desc").Limit(limit).Find(&feeds).Error
	return feeds, err
}
