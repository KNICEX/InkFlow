package dao

import "time"

type Ink struct {
	Id          int64
	Title       string `gorm:"type:varchar(100)"`
	AuthorId    int64  `gorm:"index:author_id"`
	Summary     string
	ContentMeta string
	ContentHtml string
	CreatedAt   time.Time
}

type DraftInk Ink

type LiveInk Ink
