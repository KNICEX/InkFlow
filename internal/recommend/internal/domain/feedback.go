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
	FeedbackUnLike       FeedbackType = "unlike"
	FeedbackTypeFavorite FeedbackType = "favorite"
)

func (t FeedbackType) ToString() string {
	return string(t)
}
