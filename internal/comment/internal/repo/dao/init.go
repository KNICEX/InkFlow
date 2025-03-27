package dao

import (
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	if err := db.AutoMigrate(&Comment{}, &CommentLike{}, &CommentStats{}); err != nil {
		panic(err)
	}
}
