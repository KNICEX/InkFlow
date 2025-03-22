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

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type InteractiveDAO interface {
	InsertInteractive(ctx context.Context, biz string, bizId int64) error
	InsertView(ctx context.Context, biz string, bizId, uid int64) error
	InsertFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error
	InsertLike(ctx context.Context, biz string, bizId, uid int64) error
	DeleteLike(ctx context.Context, biz string, bizId, uid int64) error
	DeleteFavorite(ctx context.Context, biz string, bizId, uid int64) error
	UpdateFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error
	InsertUnlike(ctx context.Context, biz string, bizId, uid int64) error

	InsertViewBatch(ctx context.Context, biz string, bizIds, uids []int64) error
	InsertLikeBatch(ctx context.Context, biz string, bizIds, uids []int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]Interactive, error)

	FindLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLike, error)
	FindFavoriteInfo(ctx context.Context, biz string, bizId, uid int64) (UserFavorite, error)
	FindLikeBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]UserLike, error)
	FindFavoriteBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]UserFavorite, error)
	FindByFavorite(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]UserFavorite, error)

	ListViewRecord(ctx context.Context, biz string, userId int64, maxId int64, limit int) ([]UserView, error)
	ListLikeRecord(ctx context.Context, biz string, userId int64, maxId int64, limit int) ([]UserLike, error)
}

type GormInteractiveDAO struct {
	node snowflakex.Node
	db   *gorm.DB
	l    logx.Logger
}

func NewGormInteractiveDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) InteractiveDAO {
	return &GormInteractiveDAO{
		node: node,
		db:   db,
		l:    l,
	}
}

func (dao *GormInteractiveDAO) InsertInteractive(ctx context.Context, biz string, bizId int64) error {
	now := time.Now()
	err := dao.db.WithContext(ctx).Create(&Interactive{
		Id:          dao.node.NextID(),
		Biz:         biz,
		BizId:       bizId,
		ViewCnt:     0,
		LikeCnt:     0,
		FavoriteCnt: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).Error
	err, dup := gormx.CheckDuplicateErr(err)
	if dup {
		return nil
	}
	return err
}

func (dao *GormInteractiveDAO) InsertView(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now()
	err := dao.db.WithContext(ctx).Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizId).
		UpdateColumns(map[string]any{
			"view_cnt":   gorm.Expr("interactive.view_cnt + 1"),
			"updated_at": now,
		}).Error
	if err != nil {
		dao.l.WithCtx(ctx).Error("InsertView incr view_cnt error", logx.Error(err),
			logx.Int64("userId", uid),
			logx.String("biz", biz),
			logx.Int64("bizId", bizId))
	}

	if uid == 0 {
		return nil
	}
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "biz"},
			{Name: "biz_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"updated_at": now,
		}),
	}).Create(&UserView{
		Id:        dao.node.NextID(),
		UserId:    uid,
		Biz:       biz,
		BizId:     bizId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}

func (dao *GormInteractiveDAO) InsertLike(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now()
	err := dao.db.WithContext(ctx).Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizId).
		UpdateColumns(map[string]any{
			"like_cnt":   gorm.Expr("interactive.like_cnt + 1"),
			"updated_at": now,
		}).Error
	if err != nil {
		dao.l.WithCtx(ctx).Error("InsertLike incr like_cnt error", logx.Error(err),
			logx.Int64("userId", uid),
			logx.String("biz", biz),
			logx.Int64("bizId", bizId))
	}

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "biz"},
			{Name: "biz_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"updated_at": now,
		}),
	}).Create(&UserLike{
		Id:        dao.node.NextID(),
		UserId:    uid,
		Biz:       biz,
		BizId:     bizId,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error
}

func (dao *GormInteractiveDAO) InsertFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error {
	now := time.Now()
	err := dao.db.WithContext(ctx).Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizId).
		UpdateColumns(map[string]any{
			"favorite_cnt": gorm.Expr("favorite_cnt + 1"),
			"updated_at":   now,
		}).Error
	if err != nil {
		dao.l.WithCtx(ctx).Error("InsertLike incr favorite_cnt error", logx.Error(err),
			logx.Int64("userId", uid),
			logx.String("biz", biz),
			logx.Int64("bizId", bizId))
	}

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "biz"},
			{Name: "favorite_id"},
			{Name: "biz_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"updated_at": now,
		}),
	}).Create(&UserFavorite{
		Id:         dao.node.NextID(),
		UserId:     uid,
		Biz:        biz,
		FavoriteId: fid,
		BizId:      bizId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Error
}

func (dao *GormInteractiveDAO) DeleteLike(ctx context.Context, biz string, bizId, uid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, uid).Delete(&UserLike{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		return tx.Where("biz = ? AND biz_id = ?", biz, bizId).Updates(map[string]any{
			"like_cnt":   gorm.Expr("interactive.like_cnt - 1"),
			"updated_at": time.Now(),
		}).Error
	})
}

func (dao *GormInteractiveDAO) DeleteFavorite(ctx context.Context, biz string, bizId, uid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, uid).Delete(&UserFavorite{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		return tx.Where("biz = ? AND biz_id = ?", biz, bizId).Updates(map[string]any{
			"favorite_cnt": gorm.Expr("favorite_cnt - 1"),
			"updated_at":   time.Now(),
		}).Error
	})
}

func (dao *GormInteractiveDAO) UpdateFavorite(ctx context.Context, biz string, bizId, uid, fid int64) error {
	return dao.db.WithContext(ctx).Model(&UserFavorite{}).Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, uid).
		Updates(map[string]any{
			"favorite_id": fid,
			"updated_at":  time.Now(),
		}).Error
}

func (dao *GormInteractiveDAO) InsertUnlike(ctx context.Context, biz string, bizId, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GormInteractiveDAO) InsertViewBatch(ctx context.Context, biz string, bizIds, uids []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDao := NewGormInteractiveDAO(dao.db, dao.node, dao.l)
		for i, b := range bizIds {
			err := txDao.InsertView(ctx, biz, b, uids[i])
			if err != nil {
				// 记录一下就ok
				dao.l.WithCtx(ctx).Error("InsertViewBatch error", logx.Error(err))
			}
		}
		return nil
	})
}

func (dao *GormInteractiveDAO) InsertLikeBatch(ctx context.Context, biz string, bizIds, uids []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDao := NewGormInteractiveDAO(dao.db, dao.node, dao.l)
		for i, b := range bizIds {
			err := txDao.InsertLike(ctx, biz, b, uids[i])
			if err != nil {
				// TODO 这里出错可以考虑回滚
				dao.l.WithCtx(ctx).Error("InsertLikeBatch error", logx.Error(err))
			}
		}
		return nil
	})
}

func (dao *GormInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var intr Interactive
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, bizId).First(&intr).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Interactive{}, err
	}
	return intr, nil
}

func (dao *GormInteractiveDAO) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]Interactive, error) {
	var intrs []Interactive
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id IN ?", biz, bizIds).Find(&intrs).Error
	if err != nil {
		return nil, err
	}
	res := make(map[int64]Interactive)
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}

func (dao *GormInteractiveDAO) FindLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLike, error) {
	var like UserLike
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, uid).First(&like).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return UserLike{}, err
	}
	return like, nil
}

func (dao *GormInteractiveDAO) FindFavoriteInfo(ctx context.Context, biz string, bizId, uid int64) (UserFavorite, error) {
	var favorite UserFavorite
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, uid).First(&favorite).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return UserFavorite{}, err
	}
	return favorite, nil
}

func (dao *GormInteractiveDAO) FindLikeBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]UserLike, error) {
	var likes []UserLike
	err := dao.db.WithContext(ctx).Where("biz = ? AND AND user_id = ? biz_id IN ?", biz, uid, bizIds).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	res := make(map[int64]UserLike)
	for _, like := range likes {
		res[like.BizId] = like
	}
	return res, nil
}

func (dao *GormInteractiveDAO) FindFavoriteBatch(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]UserFavorite, error) {
	var favorites []UserFavorite
	err := dao.db.WithContext(ctx).Where("biz = ? AND user_id = ? AND biz_id IN ?", biz, uid, bizIds).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	res := make(map[int64]UserFavorite)
	for _, favorite := range favorites {
		res[favorite.BizId] = favorite
	}
	return res, nil
}

func (dao *GormInteractiveDAO) FindByFavorite(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]UserFavorite, error) {
	var favorites []UserFavorite
	var tx *gorm.DB
	if maxId == 0 {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ? AND favorite_id = ?", uid, biz, fid)
	} else {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ? AND favorite_id = ? AND id < ?", uid, biz, fid, maxId)
	}
	err := tx.Order("id DESC").Limit(limit).Find(&favorites).Error
	return favorites, err
}

func (dao *GormInteractiveDAO) ListViewRecord(ctx context.Context, biz string, userId int64, maxId int64, limit int) ([]UserView, error) {
	var records []UserView
	var tx *gorm.DB
	if maxId == 0 {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ?", userId, biz)
	} else {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ? AND id < ?", userId, biz, maxId)
	}
	err := tx.Order("id DESC").Limit(limit).Find(&records).Error
	return records, err
}

func (dao *GormInteractiveDAO) ListLikeRecord(ctx context.Context, biz string, userId int64, maxId int64, limit int) ([]UserLike, error) {
	var records []UserLike
	var tx *gorm.DB
	if maxId == 0 {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ?", userId, biz)
	} else {
		tx = dao.db.WithContext(ctx).Where("user_id = ? AND biz = ? AND id < ?", userId, biz, maxId)
	}
	err := tx.Order("id DESC").Limit(limit).Find(&records).Error
	return records, err
}

type UserView struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	Biz       string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx"`
	BizId     int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}

type UserLike struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:2"`
	Biz       string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:1"`
	BizId     int64  `gorm:"uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:3"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}

// UserFavorite TODO 考虑支持多个收藏夹
type UserFavorite struct {
	Id         int64
	UserId     int64  `gorm:"uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:2"`
	Biz        string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:1"`
	FavoriteId int64  `gorm:"index;uniqueIndex:userId_biz_id_idx"`
	BizId      int64  `gorm:"uniqueIndex:userId_biz_id_idx;index:biz_uid_biz_id_idx,priority:3"`
	CreatedAt  time.Time
	UpdatedAt  time.Time `gorm:"index"`
}

type Interactive struct {
	Id          int64
	Biz         string `gorm:"type:varchar(64);uniqueIndex:biz_type_idx"`
	BizId       int64  `gorm:"uniqueIndex:biz_type_idx"`
	ViewCnt     int64  `gorm:"index"`
	LikeCnt     int64  `gorm:"index"`
	UnlikeCnt   int64
	FavoriteCnt int64 `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
