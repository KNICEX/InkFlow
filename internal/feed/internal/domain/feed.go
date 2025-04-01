package domain

import "time"

type Feed struct {
	Id        int64
	UserId    int64
	Biz       string
	BizId     int64
	Content   any
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	BizInk = "ink"
)

type FeedInk struct {
	InkId     int64     `json:"inkId"`
	AuthorId  int64     `json:"authorId"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Abstract  string    `json:"abstract"`
	CreatedAt time.Time `json:"createdAt"`
}
