package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDraftNotFound = gorm.ErrRecordNotFound
)

type DraftDAO interface {
	Insert(ctx context.Context, d DraftInk) (int64, error)
	Update(ctx context.Context, d DraftInk) error
	Delete(ctx context.Context, id int64, authorId int64, status int) error
	UpdateStatus(ctx context.Context, inkId int64, authorId int64, status int) error
	FindById(ctx context.Context, id int64) (DraftInk, error)
	FindByIdAndAuthorId(ctx context.Context, id int64, authorId int64) (DraftInk, error)
	FindByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]DraftInk, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]DraftInk, error)
}

var _ DraftDAO = (*draftDAO)(nil)

type draftDAO struct {
	db   *gorm.DB
	node snowflakex.Node
}

func NewDraftDAO(db *gorm.DB, node snowflakex.Node) DraftDAO {
	return &draftDAO{
		db:   db,
		node: node,
	}
}

func (dao *draftDAO) Insert(ctx context.Context, d DraftInk) (int64, error) {
	now := time.Now()
	d.Id = dao.node.NextID()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.Status = InkStatusUnPublished
	dao.db.WithContext(ctx).Create(&d)
	if dao.db.Error != nil {
		return 0, dao.db.Error
	}
	return d.Id, nil
}

func (dao *draftDAO) Update(ctx context.Context, d DraftInk) error {
	d.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Model(&DraftInk{}).Where("id = ? AND author_id = ?", d.Id, d.AuthorId).Updates(d).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *draftDAO) UpdateStatus(ctx context.Context, inkId int64, authorId int64, status int) error {
	err := dao.db.WithContext(ctx).Model(&DraftInk{}).Where("id = ? AND author_id = ?", inkId, authorId).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *draftDAO) FindById(ctx context.Context, id int64) (DraftInk, error) {
	var d DraftInk
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&d).Error
	if err != nil {
		return d, err
	}
	return d, nil
}

func (dao *draftDAO) FindByIdAndAuthorId(ctx context.Context, id int64, authorId int64) (DraftInk, error) {
	var d DraftInk
	err := dao.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, authorId).First(&d).Error
	if err != nil {
		return d, err
	}
	return d, nil
}

func (dao *draftDAO) Delete(ctx context.Context, id int64, authorId int64, status int) error {
	err := dao.db.WithContext(ctx).Where("id = ? AND author_id = ? AND status = ?", id, authorId, status).Delete(&DraftInk{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *draftDAO) FindByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]DraftInk, error) {
	var drafts []DraftInk
	err := dao.db.WithContext(ctx).Where("author_id = ? AND status = ?", authorId, InkStatusUnPublished).
		Order("updated_at desc").Offset(offset).Limit(limit).Find(&drafts).Error
	if err != nil {
		return drafts, err
	}
	return drafts, nil
}

func (dao *draftDAO) FindAll(ctx context.Context, maxId int64, limit int) ([]DraftInk, error) {
	var drafts []DraftInk
	err := dao.db.WithContext(ctx).Where("id < ? AND status = ?", maxId, InkStatusUnPublished).Order("id desc").Limit(limit).Find(&drafts).Error
	if err != nil {
		return drafts, err
	}
	return drafts, nil
}
