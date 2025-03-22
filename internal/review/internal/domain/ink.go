package domain

import "time"

type Ink struct {
	Id        int64
	AuthorId  int64
	Title     string
	Cover     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewResult struct {
	Passed      bool
	Reason      string
	ReviewScore int64
	ReviewTags  []string
}
