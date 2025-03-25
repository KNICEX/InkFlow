package domain

import "time"

type Ink struct {
	Id          int64
	Author      User
	Title       string
	Content     string
	Tags        []string
	AiTags      []string
	Cover       string
	ViewCnt     int64
	LikeCnt     int64
	FavoriteCnt int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InkOrderField string

const (
	OrderTypeCreatedAt InkOrderField = "created_at"
	OrderTypeViewCnt   InkOrderField = "view_cnt"
	OrderTypeLikeCnt   InkOrderField = "like_cnt"
)

type InkOrder struct {
	Field InkOrderField
	Desc  bool
}
