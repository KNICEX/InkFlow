package dao

import "time"

type Ink struct {
	Id          int64
	Title       string `gorm:"type:varchar(100)"`
	AuthorId    int64  `gorm:"index:author_id"`
	Cover       string
	Summary     string
	CategoryId  int64 `gorm:"index:category_id"`
	ContentType int   `gorm:"type:int;index:content_type"`
	ContentMeta string
	ContentHtml string
	Tags        string
	AiTags      string
	Status      int       `gorm:"type:int;default:0;index"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time `gorm:"index"`
}

type DraftInk Ink

type LiveInk Ink

const (
	InkStatusUnKnown     int = iota
	InkStatusUnPublished     // 草稿
	InkStatusPublished       // 公开
	InkStatusPrivate         // 私密
)
