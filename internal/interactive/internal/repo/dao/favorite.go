package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateFavoriteName = gorm.ErrDuplicatedKey
)

type Favorite struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:uid_biz_name"`
	Biz       string `gorm:"uniqueIndex:uid_biz_name"`
	Name      string `gorm:"uniqueIndex:uid_biz_name"`
	Private   bool   `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FavoriteDAO interface {
	Insert(ctx context.Context, f Favorite) (int64, error)
	Delete(ctx context.Context, id, uid int64) error
	FindByUserId(ctx context.Context, biz string, uid int64) ([]Favorite, error)
	Update(ctx context.Context, f Favorite) error
}

type GormFavoriteDAO struct {
	node snowflakex.Node
	db   *gorm.DB
}

func NewGormFavoriteDAO(db *gorm.DB, node snowflakex.Node) *GormFavoriteDAO {
	return &GormFavoriteDAO{
		node: node,
		db:   db,
	}
}

func (dao *GormFavoriteDAO) Insert(ctx context.Context, f Favorite) (int64, error) {
	f.Id = dao.node.NextID()
	err := dao.db.WithContext(ctx).Create(&f).Error
	err, dup := gormx.CheckDuplicateErr(err)
	if dup {
		return 0, ErrDuplicateFavoriteName
	}
	return f.Id, err
}

func (dao *GormFavoriteDAO) Delete(ctx context.Context, id, uid int64) error {
	f := Favorite{Id: id, UserId: uid}
	return dao.db.WithContext(ctx).Delete(&f).Error
}

func (dao *GormFavoriteDAO) FindByUserId(ctx context.Context, biz string, uid int64) ([]Favorite, error) {
	var fs []Favorite
	err := dao.db.WithContext(ctx).Where("user_id = ? AND biz = ?", uid, biz).Find(&fs).Error
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (dao *GormFavoriteDAO) Update(ctx context.Context, f Favorite) error {
	return dao.db.WithContext(ctx).Model(&f).Updates(f).Error
}
