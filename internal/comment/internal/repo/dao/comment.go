package dao

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/pkg/gormx"
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
	// FindByBiz 查找最新一级评论
	FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int)
	Like(ctx context.Context, userId int64, cid int64) error
	CancelLike(ctx context.Context, userId int64, cid int64) error
}

type GormCommentDAO struct {
	db *gorm.DB
}

func NewGormCommentDAO(db *gorm.DB) CommentDAO {
	return &GormCommentDAO{
		db: db,
	}
}

func (dao *GormCommentDAO) Insert(ctx context.Context, c Comment) error {
	// TODO 是否创建CommentStatistic
	return dao.db.Create(&c).Error
}

func (dao *GormCommentDAO) FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GormCommentDAO) FindRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int) {
	//TODO implement me
	panic("implement me")
}

func (dao *GormCommentDAO) insertLike(tx *gorm.DB, userId int64, cid int64) error {
	now := time.Now()
	err := tx.Create(&CommentLike{
		CommentId: cid,
		UserId:    userId,
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
