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
	Passed      bool     `json:"passed"`
	Reason      string   `json:"reason"`
	ReviewScore int64    `json:"reviewScore"`
	ReviewTags  []string `json:"reviewTags"`
}
