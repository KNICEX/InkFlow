package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateKey   = gorm.ErrDuplicatedKey
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechatOpenId(ctx context.Context, openId string) (User, error)
	UpdateById(ctx context.Context, u User) error
	FindByGithubId(ctx context.Context, id int64) (User, error)
	FindByAccountName(ctx context.Context, accountName string) (User, error)
}

type User struct {
	Id       int64          `gorm:"primaryKey"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password sql.NullString `gorm:"type:varchar(100)"`
	Account  string         `gorm:"unique"`
	Username string         `gorm:"type:varchar(60)"`
	AboutMe  string         `gorm:"type:varchar(500)"`
	Exp      int64          `gorm:"default:0"`
	Level    int            `gorm:"default:1"`
	// 逗号分隔
	Links      string            `gorm:"type:varchar(500)"`
	Birthday   sql.NullTime      `gorm:"type:date"`
	GithubId   sql.NullInt64     `gorm:"unique"`
	GithubInfo domain.GithubInfo `gorm:"serializer:json;type:json"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewGormUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicateKey
	}

	return err
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *GormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *GormUserDAO) FindByWechatOpenId(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *GormUserDAO) FindByGithubId(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("github_id = ?", id).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *GormUserDAO) UpdateById(ctx context.Context, u User) error {
	u.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Where("id = ?", u.Id).Updates(u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return err
}

func (dao *GormUserDAO) FindByAccountName(ctx context.Context, accountName string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("account_name = ?", accountName).First(&u).Error
	if err != nil {
		return u, err
	}
	return u, nil
}
