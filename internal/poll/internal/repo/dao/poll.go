package dao

import "time"

type Poll struct {
	Id        int64
	CreatorId int64
	BizId     int64
	BizType   string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PollOption struct {
	Id        int64
	Title     string
	PollId    int64
	Count     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PollRecord struct {
	Id        int64
	PollId    int64
	OptionId  int64
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
