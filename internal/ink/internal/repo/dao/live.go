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
	UpdateAiTags(ctx context.Context, inkId int64, aiTags string) error
	Delete(ctx context.Context, id int64, authorId int64, status ...int) error
	FindById(ctx context.Context, id int64, status ...int) (LiveInk, error)
	FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...int) ([]LiveInk, error)
	FindAll(ctx context.Context, maxId int64, limit int, status ...int) ([]LiveInk, error)
	FindFirstPageByAuthorIds(ctx context.Context, authorIds []int64, n int, status ...int) (map[int64][]LiveInk, error)
	FindByIds(ctx context.Context, ids []int64, status ...int) (map[int64]LiveInk, error)
	CountByAuthorId(ctx context.Context, authorId int64, status ...int) (int64, error)
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
	err := dao.db.WithContext(ctx).Model(&LiveInk{}).Where("id = ? AND author_id = ?", inkId, authorId).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *liveDAO) UpdateAiTags(ctx context.Context, inkId int64, aiTags string) error {
	err := dao.db.WithContext(ctx).Model(&LiveInk{}).Where("id = ?", inkId).
		Update("ai_tags", aiTags).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *liveDAO) FindById(ctx context.Context, id int64, status ...int) (LiveInk, error) {
	var ink LiveInk
	var err error
	if len(status) > 0 {
		err = dao.db.WithContext(ctx).Where("id = ? AND status IN ?", id, status).First(&ink).Error
	} else {
		err = dao.db.WithContext(ctx).Where("id = ?", id).First(&ink).Error
	}
	if err != nil {
		return LiveInk{}, err
	}
	return ink, nil
}

func (dao *liveDAO) Delete(ctx context.Context, id int64, authorId int64, status ...int) error {
	var err error
	if len(status) > 0 {
		err = dao.db.WithContext(ctx).Where("id = ? AND author_id = ? AND status IN ?", id, authorId, status).Delete(&LiveInk{}).Error
	} else {
		err = dao.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, authorId).Delete(&LiveInk{}).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func (dao *liveDAO) FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...int) ([]LiveInk, error) {
	var inks []LiveInk
	var err error
	if len(status) > 0 {
		err = dao.db.WithContext(ctx).Where("author_id = ? AND status IN ?", authorId, status).
			Order("updated_at DESC").Offset(offset).Limit(limit).Find(&inks).Error
	} else {
		err = dao.db.WithContext(ctx).Where("author_id = ?", authorId).
			Order("updated_at DESC").Offset(offset).Limit(limit).Find(&inks).Error
	}
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindAll(ctx context.Context, maxId int64, limit int, status ...int) ([]LiveInk, error) {
	var inks []LiveInk

	var tx = dao.db
	if maxId != 0 {
		tx = tx.Where("id < ?", maxId)
	}
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}

	err := tx.WithContext(ctx).Order("id DESC").Limit(limit).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (dao *liveDAO) FindByIds(ctx context.Context, ids []int64, status ...int) (map[int64]LiveInk, error) {
	var inks []LiveInk
	var err error
	if len(status) > 0 {
		err = dao.db.WithContext(ctx).Where("id IN ? AND status IN ?", ids, status).Find(&inks).Error
	} else {
		err = dao.db.WithContext(ctx).Where("id IN ?", ids).Find(&inks).Error
	}
	if err != nil {
		return nil, err
	}
	idMap := make(map[int64]LiveInk)
	for _, ink := range inks {
		idMap[ink.Id] = ink
	}
	return idMap, nil
}

func (dao *liveDAO) FindFirstPageByAuthorIds(ctx context.Context, authorIds []int64, n int, status ...int) (map[int64][]LiveInk, error) {
	var inks []LiveInk

	subQuery := dao.db.WithContext(ctx).Model(&LiveInk{}).Where("author_id IN ?", authorIds)
	if len(status) > 0 {
		subQuery = subQuery.Where("status IN ?", status)
	}
	subQuery = subQuery.Select("*, ROW_NUMBER() OVER (PARTITION BY author_id ORDER BY updated_at DESC) AS row_num")

	err := dao.db.WithContext(ctx).Table("(?) as t", subQuery).
		Where("t.row_num <= ?", n).Find(&inks).Error
	if err != nil {
		return nil, err
	}
	inkMap := make(map[int64][]LiveInk)
	for _, ink := range inks {
		inkMap[ink.AuthorId] = append(inkMap[ink.AuthorId], ink)
	}
	return inkMap, nil
}

func (dao *liveDAO) CountByAuthorId(ctx context.Context, authorId int64, status ...int) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&LiveInk{}).Where("author_id = ? AND status IN ?", authorId, status).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
