package dao

import "time"

type Tag struct {
	Id        int64
	Name      string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
