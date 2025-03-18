package domain

import "time"

type SubscribeRelation struct {
	SubscriberUid int64
	SubscribedUid int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
