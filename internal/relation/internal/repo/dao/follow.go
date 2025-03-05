package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"gorm.io/gorm"
	"time"
)

type UserFollow struct {
	Id       int64
	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee,index:followee"`

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
//	Followings int64
//
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}

type FollowRelationDAO interface {
	FollowList(ctx context.Context, uid, maxId int64, limit int) ([]UserFollow, error)
	CreateFollowRelation(ctx context.Context, c UserFollow) error
	CancelFollow(ctx context.Context, c UserFollow) error
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
	//FollowStatistic(ctx context.Context, uid int64) (FollowStatistic, error)
}

type GormFollowRelationDAO struct {
	db *gorm.DB
	l  logx.Logger
}

func NewGormFollowRelationDAO(db *gorm.DB, l logx.Logger) FollowRelationDAO {
	return &GormFollowRelationDAO{
		db: db,
		l:  l,
	}
}

func (dao *GormFollowRelationDAO) FollowList(ctx context.Context, uid, maxId int64, limit int) ([]UserFollow, error) {
	var res []UserFollow
	var tx *gorm.DB
	if maxId == 0 {
		tx = dao.db.WithContext(ctx).Where("follower_id = ?", uid)
	} else {
		tx = dao.db.WithContext(ctx).Where("follower_id = ? AND id < ?", uid, maxId)
	}
	if err := tx.Order("id DESC").Limit(limit).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) CreateFollowRelation(ctx context.Context, c UserFollow) error {
	err := dao.db.WithContext(ctx).Create(&c).Error
	err, ok := gormx.CheckDuplicateErr(err)
	if ok {
		// 已经关注过
		dao.l.Warn("连续follow", logx.Int64("UserId", c.Follower))
		return nil
	}
	return err
}
func (dao *GormFollowRelationDAO) CancelFollow(ctx context.Context, c UserFollow) error {
	return dao.db.WithContext(ctx).Delete(&c).Error
}
func (dao *GormFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	if err := dao.db.WithContext(ctx).Where("followee_id = ?", uid).Count(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	if err := dao.db.WithContext(ctx).Where("follower_id = ?", uid).Count(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}

//func (dao *GormFollowRelationDAO) FollowStatistic(ctx context.Context, uid int64) (FollowStatistic, error) {
//	panic("implement me")
//}
