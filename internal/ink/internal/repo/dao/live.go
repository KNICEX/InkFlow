package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
	ErrLiveInkNotFound = gorm.ErrRecordNotFound
)

type LiveDAO interface {
	Upsert(ctx context.Context, d LiveInk) (int64, error)
	UpdateStatus(ctx context.Context, inkId int64, authorId int64, status int) error
	Delete(ctx context.Context, id int64, authorId int64, status []int) error
	FindById(ctx context.Context, id int64) (LiveInk, error)
	FindByIdAndStatus(ctx context.Context, id int64, status int) (LiveInk, error)
	FindByAuthorIdAndMaxId(ctx context.Context, authorId int64, maxId int64, limit int) ([]LiveInk, error)
	FindByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]LiveInk, error)
	FindByAuthorIdAndStatus(ctx context.Context, authorId int64, status int, offset, limit int) ([]LiveInk, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]LiveInk, error)
	FindAllByStatus(ctx context.Context, status int, maxId int64, limit int) ([]LiveInk, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]LiveInk, error)
}

var _ LiveDAO = (*liveDAO)(nil)

type liveDAO struct {
	db *gorm.DB
}

func NewLiveDAO(db *gorm.DB) LiveDAO {
	return &liveDAO{
		db: db,
	}
}

func (dao *liveDAO) Upsert(ctx context.Context, d LiveInk) (int64, error) {
	now := time.Now()
	d.UpdatedAt = now
	d.CreatedAt = now
	return d.Id, dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"title",
			"cover",
			"summary",
			"category_id",
			"status",
			"content_html",
			"content_meta",
			"tags",
			"ai_tags",
			"updated_at",
		}),
	}).Create(&d).Error
}

func (dao *liveDAO) UpdateStatus(ctx context.Context, inkId int64, authorId int64, status int) error {
	err := dao.db.WithContext(ctx).Model(&LiveInk{}).Where("id = ? and author_id = ?", inkId, authorId).
		Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *liveDAO) FindById(ctx context.Context, id int64) (LiveInk, error) {
	var ink LiveInk
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&ink).Error
	if err != nil {
		return LiveInk{}, err
	}
	return ink, nil
}

func (dao *liveDAO) Delete(ctx context.Context, id int64, authorId int64, status []int) error {
	err := dao.db.WithContext(ctx).Where("id = ? and author_id = ? AND status in ?", id, authorId, status).Delete(&LiveInk{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *liveDAO) FindByIdAndStatus(ctx context.Context, id int64, status int) (LiveInk, error) {
	var ink LiveInk
	err := dao.db.WithContext(ctx).Where("id = ? and status = ?", id, status).First(&ink).Error
	if err != nil {
		return LiveInk{}, err
	}
	return ink, nil
}

func (dao *liveDAO) FindByAuthorIdAndMaxId(ctx context.Context, authorId int64, maxId int64, limit int) ([]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("author_id = ? and id < ?", authorId, maxId).
		Order("id desc").Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("author_id = ?", authorId).
		Order("updated_at desc").Offset(offset).Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindByAuthorIdAndStatus(ctx context.Context, authorId int64, status int, offset, limit int) ([]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("author_id = ? and status = ?", authorId, status).
		Order("updated_at desc").Offset(offset).Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindAll(ctx context.Context, maxId int64, limit int) ([]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("id < ?", maxId).
		Order("id desc").Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindAllByStatus(ctx context.Context, status int, maxId int64, limit int) ([]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("status = ? and id < ?", status, maxId).
		Order("updated_at desc").Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindByIds(ctx context.Context, ids []int64) (map[int64]LiveInk, error) {
	var inks []LiveInk
	err := dao.db.WithContext(ctx).Where("id in ?", ids).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	idMap := make(map[int64]LiveInk)
	for _, ink := range inks {
		idMap[ink.Id] = ink
	}
	return idMap, nil
}
