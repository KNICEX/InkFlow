package domain

import "time"

type FeedEvent struct {
	Id        int64
	Type      string
	CreatedAt time.Time
}
