package domain

import "time"

type Ink struct {
	Id          int64
	Title       string
	AuthorId    int64
	CategoryId  int64
	ContentType int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
