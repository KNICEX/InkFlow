package event

import "time"

type FollowEvt struct {
	FollowerId int64
	FolloweeId int64
	CreatedAt  time.Time
}
