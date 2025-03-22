package domain

import "time"

type Feedback struct {
	FeedbackType FeedbackType
	UserId       int64
	InkId        int64
	CreatedAt    time.Time
}

type FeedbackType string

const (
	FeedbackTypeView     FeedbackType = "view"
	FeedbackTypeViewLong FeedbackType = "view_long"
	FeedbackTypeLike     FeedbackType = "like"
	FeedbackTypeFavorite FeedbackType = "favorite"
)
