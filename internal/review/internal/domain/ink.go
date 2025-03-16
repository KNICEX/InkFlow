package domain

import "time"

type Ink struct {
	Id        int64
	AuthorId  int64
	Title     string
	Cover     string
	Content   string
	Status    int
	Tag       []string
	AiTag     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
