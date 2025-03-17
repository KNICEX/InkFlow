package dao

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

type UserFollow struct {
	Id         int64
	FollowerId int64 `gorm:"uniqueIndex:follower_followee,index:follower"`
	FolloweeId int64 `gorm:"uniqueIndex:follower_followee,index:followee"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

//type Block struct {
//	Id        int64
//	Blocker   int64
//	Blocked   int64
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}

//type FollowStatistic struct {
//	Id         int64
//	UserId     int64 `gorm:"unique"`
//	Followers  int64
//	Following int64
//
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}

type FollowRelationDAO interface {
	FollowList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error)
	FollowerList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error)
	CreateFollowRelation(ctx context.Context, c UserFollow) error
	CancelFollow(ctx context.Context, c UserFollow) error
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
	Followed(ctx context.Context, uid, followeeId int64) (bool, error)
	//FollowStatistic(ctx context.Context, uid int64) (FollowStatistic, error)
}

type GormFollowRelationDAO struct {
	db   *gorm.DB
	node snowflakex.Node
	l    logx.Logger
}

func NewGormFollowRelationDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) FollowRelationDAO {
	return &GormFollowRelationDAO{
		db:   db,
		node: node,
		l:    l,
	}
}

func (dao *GormFollowRelationDAO) FollowList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error) {
	var res []UserFollow
	tx := dao.db
	var err error
	if maxId == 0 {
		err = tx.WithContext(ctx).Where("follower_id = ?", uid).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.WithContext(ctx).Where("follower_id = ? AND id < ?", uid, maxId).Order("id DESC").Limit(limit).Find(&res).Error
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) FollowerList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error) {
	var res []UserFollow
	tx := dao.db
	var err error
	if maxId == 0 {
		err = tx.WithContext(ctx).Where("followee_id = ?", uid).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.WithContext(ctx).Where("followee_id = ? AND id < ?", uid, maxId).Order("id DESC").Limit(limit).Find(&res).Error
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) CreateFollowRelation(ctx context.Context, c UserFollow) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Id = dao.node.NextID()
	err := dao.db.WithContext(ctx).Create(&c).Error
	err, ok := gormx.CheckDuplicateErr(err)
	if ok {
		// 已经关注过
		dao.l.Warn("连续follow", logx.Int64("UserId", c.FollowerId))
		return nil
	}
	return err
}
func (dao *GormFollowRelationDAO) CancelFollow(ctx context.Context, c UserFollow) error {
	return dao.db.WithContext(ctx).Delete(&c).Error
}
func (dao *GormFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	if err := dao.db.WithContext(ctx).Model(&UserFollow{}).Where("followee_id = ?", uid).Count(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	if err := dao.db.WithContext(ctx).Model(&UserFollow{}).Where("follower_id = ?", uid).Count(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) Followed(ctx context.Context, uid, followeeId int64) (bool, error) {
	var res UserFollow
	err := dao.db.WithContext(ctx).Where("follower_id = ? AND followee_id = ?", uid, followeeId).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//func (dao *GormFollowRelationDAO) FollowStatistic(ctx context.Context, uid int64) (FollowStatistic, error) {
//	panic("implement me")
//}
