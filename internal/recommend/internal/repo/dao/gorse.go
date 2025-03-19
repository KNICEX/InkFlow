package dao

import (
	"context"
	"time"
)

type Item struct {
	ItemId     string
	IsHidden   bool
	Categories []string
	Timestamp  time.Time
	Labels     []string
	Comment    string
}

type User struct {
	UserId string   // user ID
	Labels []string // labels describing the user
}

type Feedback struct {
	FeedbackType string    // feedback type
	UserId       string    // user id
	ItemId       string    // item id
	Timestamp    time.Time // feedback timestamp
}

type GorseDAO interface {
	InsertItem(ctx context.Context, item Item) error
	HiddenItem(ctx context.Context, itemId string) error
	ReShowItem(ctx context.Context, itemId string) error

	InsertFeedback(ctx context.Context, feedback Feedback) error
	InsertUser(ctx context.Context, user User) error
}
