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

type Comment struct {
	Id    int64  `gorm:"primaryKey"`
	Biz   string `gorm:"index:biz_type_id"`
	BizId int64  `gorm:"index:biz_type_id"`

	CommentatorId int64
	IsAuthor      bool
	Content       string
	Images        string

	// 根评论id
	RootId int64 `gorm:"index"`
	// 父评论id
	ParentId int64 `gorm:"index"`

	CreatedAt time.Time
}

type CommentLike struct {
	Id        int64
	CommentId int64 `gorm:"uniqueIndex:comment_user_id"`
	UserId    int64 `gorm:"uniqueIndex:comment_user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentStats struct {
	Id         int64
	CommentId  int64 `gorm:"unique"`
	LikeCount  int64
	ReplyCount int64
	Heat       int64 `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CommentDAO interface {
	Insert(ctx context.Context, c Comment) (int64, error)

	CommentCnt(ctx context.Context, biz string, bizId int64) (int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteByBiz(ctx context.Context, biz string, bizId int64) error
	// FindByBiz 查找最新一级评论
	FindByBiz(ctx context.Context, biz string, bizId int64, maxId int64, limit int) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, maxId int64, limit int) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int) ([]Comment, error)
	Like(ctx context.Context, userId int64, cid int64) error
	CancelLike(ctx context.Context, userId int64, cid int64) error
	FindByIds(ctx context.Context, ids []int64) (map[int64]Comment, error)
	FindById(ctx context.Context, id int64) (Comment, error)
	FindAuthorReplyIn(ctx context.Context, ids []int64) (map[int64][]Comment, error)
	FindStats(ctx context.Context, ids []int64) (map[int64]CommentStats, error)
	Liked(ctx context.Context, uid int64, cids []int64) (map[int64]bool, error)
	ReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error)
}

type GormCommentDAO struct {
	db   *gorm.DB
	node snowflakex.Node
	l    logx.Logger
}

func NewGormCommentDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) CommentDAO {
	return &GormCommentDAO{
		db:   db,
		node: node,
		l:    l,
	}
}

func (dao *GormCommentDAO) Insert(ctx context.Context, c Comment) (int64, error) {
	// TODO 是否创建CommentStatistic
	id := dao.node.NextID()
	return id, dao.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		c.Id = id
		c.CreatedAt = now
		err := tx.WithContext(ctx).Create(&c).Error
		if err != nil {
			return err
		}
		if c.ParentId == 0 {
			return nil
		}

		// 更新根评论的回复数
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "comment_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"reply_count": gorm.Expr("comment_stats.reply_count + 1"),
				"updated_at":  now,
			}),
		}).Create(&CommentStats{
			Id:         dao.node.NextID(),
			CommentId:  c.RootId,
			LikeCount:  0,
			ReplyCount: 1,
			Heat:       1,
			CreatedAt:  now,
			UpdatedAt:  now,
		}).Error

	})
}

func (dao *GormCommentDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		c := Comment{}
		err := tx.WithContext(ctx).Where("id = ?", id).First(&c).Error
		if err != nil {
			return err
		}

		if c.RootId != 0 {
			// 扣除根评论的回复数
			err = tx.WithContext(ctx).Model(&CommentStats{}).Where("comment_id = ?", c.RootId).Updates(map[string]any{
				"reply_count": gorm.Expr("reply_count - 1"),
				"updated_at":  time.Now(),
			}).Error
			if err != nil {
				return err
			}
		}

		err = tx.WithContext(ctx).Where("id = ?", id).Delete(&Comment{}).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Where("parent_id = ? OR root_id = ?", id, id).Delete(&Comment{}).Error
	})
}

func (dao *GormCommentDAO) DeleteByBiz(ctx context.Context, biz string, bizId int64) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, bizId).Delete(&Comment{}).Error
		if err != nil {
			return err
		}
		// TODO 考虑清除统计数据
		return nil
	})
}

func (dao *GormCommentDAO) CommentCnt(ctx context.Context, biz string, bizId int64) (int64, error) {
	var cnt int64
	// 只统计直接评论数
	err := dao.db.WithContext(ctx).Model(&Comment{}).Where("biz = ? AND biz_id = ? AND parent_id = 0", biz, bizId).Count(&cnt).Error
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
		err = tx.Where("biz = ? AND biz_id = ? AND parent_id = 0", biz, bizId).Order("id DESC").Limit(limit).Find(&res).Error
	} else {
		err = tx.Where("biz = ? AND biz_id = ? AND parent_id = 0 AND id < ?", biz, bizId, maxId).Order("id DESC").Limit(limit).Find(&res).Error
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

func (dao *GormCommentDAO) Like(ctx context.Context, userId int64, cid int64) error {
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		err := tx.WithContext(ctx).Create(&CommentLike{
			Id:        dao.node.NextID(),
			CommentId: cid,
			UserId:    userId,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
		err, _ = gormx.CheckDuplicateErr(err)
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "comment_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("comment_stats.like_count + 1"),
				"updated_at": now,
			}),
		}).Create(&CommentStats{
			Id:        dao.node.NextID(),
			CommentId: cid,
			LikeCount: 1,
			Heat:      1,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 重复插入，说明已经点赞了
		return nil
	}
	return err
}

func (dao *GormCommentDAO) CancelLike(ctx context.Context, userId int64, cid int64) error {
	err := dao.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userId, cid).Delete(&CommentLike{}).Error
	if err != nil {
		return nil
	}
	err = dao.db.WithContext(ctx).Model(&CommentStats{}).Where("comment_id = ?", cid).
		Update("like_count", gorm.Expr("like_count - 1")).Error
	if err != nil {
		// 这里失败了不影响业务逻辑，记录日志即可
		dao.l.WithCtx(ctx).Error("decr comment like count error", logx.Error(err),
			logx.Int64("commentId", cid))
	}
	return nil
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

func (dao *GormCommentDAO) FindById(ctx context.Context, id int64) (Comment, error) {
	var res Comment
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	if err != nil {
		return Comment{}, err
	}
	return res, nil
}

func (dao *GormCommentDAO) FindAuthorReplyIn(ctx context.Context, ids []int64) (map[int64][]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).Where("root_id IN ? AND is_author = true", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[int64][]Comment, len(res))
	for _, c := range res {
		resMap[c.RootId] = append(resMap[c.RootId], c)
	}
	return resMap, nil
}

func (dao *GormCommentDAO) FindStats(ctx context.Context, ids []int64) (map[int64]CommentStats, error) {
	var res []CommentStats
	err := dao.db.WithContext(ctx).Where("comment_id IN ?", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[int64]CommentStats, len(res))
	for _, c := range res {
		resMap[c.CommentId] = c
	}
	return resMap, nil
}

func (dao *GormCommentDAO) Liked(ctx context.Context, uid int64, cids []int64) (map[int64]bool, error) {
	var res []CommentLike
	err := dao.db.WithContext(ctx).Where("user_id = ? AND comment_id IN ?", uid, cids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[int64]bool, len(res))
	for _, c := range res {
		resMap[c.CommentId] = true
	}
	return resMap, nil
}

func (dao *GormCommentDAO) ReplyCount(ctx context.Context, biz string, bizIds []int64) (map[int64]int64, error) {
	type ReplyCount struct {
		BizId      int64 `gorm:"index:biz_type_id"`
		ReplyCount int64
	}
	var res []ReplyCount
	err := dao.db.WithContext(ctx).Model(&Comment{}).Select("biz_id, count(*) as reply_count").
		Where("biz = ? AND biz_id IN ? AND parent_id = 0", biz, bizIds).
		Group("biz_id").
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	replyCountMap := make(map[int64]int64, len(res))
	for _, item := range res {
		replyCountMap[item.BizId] = item.ReplyCount
	}
	return replyCountMap, nil
}
