package domain

import "time"

type Ink struct {
	Id        int64
	AuthorId  int64
	Category  int
	Tags      []string
	CreatedAt time.Time
}
