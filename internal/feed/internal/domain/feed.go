package domain

import "time"

type Feed struct {
	Id        int64
	UserId    int64
	FeedType  FeedType
	FeedId    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FeedType = string

const (
	FeedTypeInk FeedType = "ink"
)

type FeedInk struct {
	InkId     int64     `json:"inkId"`
	AuthorId  int64     `json:"authorId"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Abstract  string    `json:"abstract"`
	CreatedAt time.Time `json:"createdAt"`
}
