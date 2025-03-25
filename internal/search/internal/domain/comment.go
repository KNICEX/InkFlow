package domain

import "time"

type Comment struct {
	Id          int64
	Biz         string
	BizId       int64
	RootId      int64
	ParentId    int64
	Content     string
	Commentator User
	CreatedAt   time.Time
}

const BizInk = "ink"
