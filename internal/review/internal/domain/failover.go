package domain

import "time"

type FailReview struct {
	Id        int64
	Type      ReviewType
	Event     any
	Error     error
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewType string

const (
	ReviewTypeInk ReviewType = "ink"
)
