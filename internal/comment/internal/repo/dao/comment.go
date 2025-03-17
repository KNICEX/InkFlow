package dao

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Comment struct {
	Id      int64 `gorm:"primaryKey"`
	UserId  int64
	Biz     string `gorm:"index:biz_type_id"`
	BizId   int64  `gorm:"index:biz_type_id"`
	Content string
	// 根评论id
	RootId int64 `gorm:"index,default:-1"`
	// 父评论id
	ParentId    int64 `gorm:"index,default:-1"`
	ReplyUserId int64
	Status      CommentStatus

	CreatedAt time.Time
	// 基本不允许修改
	UpdatedAt time.Time
}
type CommentStatus uint8

const (
	CommentStatusUnknown CommentStatus = iota
	CommentStatusReviewing
	CommentStatusPassed
)

type CommentLike struct {
	Id        int64
	CommentId int64 `gorm:"uniqueIndex:comment_user_id"`
	UserId    int64 `gorm:"uniqueIndex:comment_user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type CommentLikeStatus uint8

const (
	CommentLikeStatusUnknown CommentLikeStatus = iota
	CommentLikeStatusLiked
	CommentLikeStatusCanceled
)

type CommentStatistic struct {
	Id         int64
	CommentId  int64 `gorm:"unique"`
	LikeCount  int64
	ReplyCount int64
	Heat       int64 `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CommentDAO interface {
	Insert(ctx context.Context, c Comment) error

	CommentCnt(ctx context.Context, biz string, bizId int64) (int64, error)
	// FindByBiz 查找最新一级评论
	FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int) ([]Comment, error)
	Like(ctx context.Context, userId int64, cid int64) error
	CancelLike(ctx context.Context, userId int64, cid int64) error
	FindByIds(ctx context.Context, ids []int64) (map[int64]Comment, error)
}

type GormCommentDAO struct {
	db   *gorm.DB
	node snowflakex.Node
}

func NewGormCommentDAO(db *gorm.DB) CommentDAO {
	return &GormCommentDAO{
		db: db,
	}
}

func (dao *GormCommentDAO) Insert(ctx context.Context, c Comment) error {
	// TODO 是否创建CommentStatistic
	c.Id = dao.node.NextID()
	return dao.db.WithContext(ctx).Create(&c).Error
}

func (dao *GormCommentDAO) CommentCnt(ctx context.Context, biz string, bizId int64) (int64, error) {
	var cnt int64
	err := dao.db.WithContext(ctx).Model(&Comment{}).Where("biz = ? AND biz_id = ?", biz, bizId).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (dao *GormCommentDAO) FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]Comment, error) {
	var res []Comment
	var err error
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		err = tx.Where("biz = ? AND biz_id = ? AND parent_id = -1", biz, bizId).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.Where("biz = ? AND biz_id = ? AND parent_id = -1 AND id < ?", biz, bizId, maxId).Order("id DESC").Limit(limit).Find(&res).Error
	}
	return res, err
}

func (dao *GormCommentDAO) FindRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int) ([]Comment, error) {
	var res []Comment
	var err error
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		err = tx.Where("parent_id = ?", pid).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.Where("parent_id = ? AND id < ?", pid, maxId).Order("id DESC").Limit(limit).Find(&res).Error
	}
	return res, err
}

func (dao *GormCommentDAO) FindRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int) ([]Comment, error) {
	var res []Comment
	var err error
	tx := dao.db.WithContext(ctx)
	if maxId == 0 {
		err = tx.Where("root_id = ?", rid).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.Where("root_id = ? AND id < ?", rid, maxId).Order("id DESC").Limit(limit).Find(&res).Error
	}
	return res, err
}

func (dao *GormCommentDAO) insertLike(tx *gorm.DB, userId int64, cid int64) error {
	now := time.Now()
	err := tx.Create(&CommentLike{
		Id:        dao.node.NextID(),
		CommentId: cid,
		UserId:    userId,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error
	er, _ := gormx.CheckDuplicateErr(err)
	if err != nil {
		return er
	}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "comment_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"updated_at": now,
			"like_count": gorm.Expr("comment_statistics.like_count + 1"),
		}),
	}).Create(&CommentStatistic{
		CommentId: cid,
		LikeCount: 1,
		Heat:      1,
		UpdatedAt: now,
		CreatedAt: now,
	}).Error
}

func (dao *GormCommentDAO) Like(ctx context.Context, userId int64, cid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ? AND comment_id = ?", userId, cid).First(&CommentLike{}).Error
		switch {
		case err == nil:
			// 存在点赞记录, 不做任何处理
			return nil
		case errors.Is(err, gorm.ErrRecordNotFound):
			return dao.insertLike(tx, userId, cid)
		default:
			return err
		}
	})
}

func (dao *GormCommentDAO) deleteLike(tx *gorm.DB, userId int64, cid int64) error {
	err := tx.Where("user_id = ? AND comment_id = ?", userId, cid).Delete(&CommentLike{}).Error
	if err != nil {
		return err
	}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "comment_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"updated_at": gorm.Expr("comment_statistics.updated_at"),
			"like_count": gorm.Expr("comment_statistics.like_count - 1"),
		}),
	}).Create(&CommentStatistic{
		CommentId: cid,
		LikeCount: -1,
		Heat:      -1,
		UpdatedAt: time.Now(),
	}).Error
}

func (dao *GormCommentDAO) CancelLike(ctx context.Context, userId int64, cid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ? AND comment_id = ?", userId, cid).First(&CommentLike{}).Error
		switch {
		case err == nil:
			return dao.deleteLike(tx, userId, cid)
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil
		default:
			return err
		}
	})
}

func (dao *GormCommentDAO) FindByIds(ctx context.Context, ids []int64) (map[int64]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[int64]Comment, len(res))
	for _, c := range res {
		resMap[c.Id] = c
	}
	return resMap, nil
}
