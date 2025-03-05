package domain

import "time"

type Poll struct {
	Id        int64
	CreatorId int64
	BizType   string
	BizId     int64
	Title     string
	Options   []PollOption
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PollOption struct {
	Id        int64
	Title     string
	Count     int64
	Polled    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
