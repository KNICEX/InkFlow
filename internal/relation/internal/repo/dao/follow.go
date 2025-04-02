package dao

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserFollow struct {
	Id         int64
	FollowerId int64 `gorm:"uniqueIndex:follower_followee,index:follower"`
	FolloweeId int64 `gorm:"uniqueIndex:follower_followee,index:followee"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

type FollowStats struct {
	Id        int64
	UserId    int64 `gorm:"unique"`
	Followers int64
	Following int64
	CreatedAt time.Time
}

var (
	ErrFollowExist = gorm.ErrDuplicatedKey
)

type FollowRelationDAO interface {
	FollowList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error)
	FollowerList(ctx context.Context, uid int64, maxId int64, limit int) ([]UserFollow, error)
	CreateFollowRelation(ctx context.Context, c UserFollow) error
	CancelFollow(ctx context.Context, c UserFollow) error
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
	FindFollowStats(ctx context.Context, uid int64) (FollowStats, error)
	FindFollowStatsBatch(ctx context.Context, uids []int64) (map[int64]FollowStats, error)
	Followed(ctx context.Context, uid, followeeId int64) (bool, error)
	FollowedBatch(ctx context.Context, uid int64, followeeIds []int64) (map[int64]bool, error)
	//FollowStats(ctx context.Context, uid int64) (FollowStats, error)
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
	return dao.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		c.CreatedAt = now
		c.UpdatedAt = now
		c.Id = dao.node.NextID()
		// create follow relation
		err := dao.db.WithContext(ctx).Create(&c).Error
		err, dup := gormx.CheckDuplicateErr(err)
		if dup {
			// 已经关注过
			dao.l.Warn("连续follow", logx.Int64("Uid", c.FollowerId))
			return nil
		}
		if err != nil {
			return err
		}

		// incr followee`s followers count
		if err = dao.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"followers": gorm.Expr("follow_stats.followers + 1"),
			}),
		}).Create(&FollowStats{
			Id:        dao.node.NextID(),
			UserId:    c.FolloweeId,
			Followers: 1,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}

		// incr follower`s following count
		if err = dao.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"following": gorm.Expr("follow_stats.following + 1"),
			}),
		}).Create(&FollowStats{
			Id:        dao.node.NextID(),
			UserId:    c.FollowerId,
			Following: 1,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}
func (dao *GormFollowRelationDAO) CancelFollow(ctx context.Context, c UserFollow) error {
	//return dao.db.WithContext(ctx).Delete(&c).Error
	return dao.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		c.UpdatedAt = now
		// delete follow relation
		err := dao.db.WithContext(ctx).Where("follower_id = ? AND followee_id = ?", c.FollowerId, c.FolloweeId).Delete(&c).Error
		if err != nil {
			return err
		}

		// decr followee`s followers count
		if err = dao.db.WithContext(ctx).Model(&FollowStats{}).Where("user_id = ?", c.FolloweeId).Update("followers", gorm.Expr("follow_stats.followers - 1")).Error; err != nil {
			return err
		}

		// decr follower`s following count
		if err = dao.db.WithContext(ctx).Model(&FollowStats{}).Where("user_id = ?", c.FollowerId).Update("following", gorm.Expr("follow_stats.following - 1")).Error; err != nil {
			return err
		}
		return nil
	})
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

func (dao *GormFollowRelationDAO) FindFollowStats(ctx context.Context, uid int64) (FollowStats, error) {
	var res FollowStats
	err := dao.db.WithContext(ctx).Where("user_id = ?", uid).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return FollowStats{
				UserId: uid,
			}, nil
		}
		return FollowStats{}, err
	}
	return res, nil
}

func (dao *GormFollowRelationDAO) FindFollowStatsBatch(ctx context.Context, uids []int64) (map[int64]FollowStats, error) {
	var res []FollowStats
	err := dao.db.WithContext(ctx).Where("user_id IN ?", uids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[int64]FollowStats)
	for _, item := range res {
		stats[item.UserId] = item
	}
	return stats, nil
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

func (dao *GormFollowRelationDAO) FollowedBatch(ctx context.Context, uid int64, followeeIds []int64) (map[int64]bool, error) {
	var res []UserFollow
	err := dao.db.WithContext(ctx).Where("follower_id = ? AND followee_id IN ?", uid, followeeIds).Find(&res).Error
	if err != nil {
		return nil, err
	}
	followed := make(map[int64]bool)
	for _, item := range res {
		followed[item.FolloweeId] = true
	}
	return followed, nil
}

//func (dao *GormFollowRelationDAO) FollowStats(ctx context.Context, uid int64) (FollowStats, error) {
//	panic("implement me")
//}
