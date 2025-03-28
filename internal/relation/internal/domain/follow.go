package domain

import "time"

type FollowRelation struct {
	FollowerId int64
	FolloweeId int64
	CreatedAt  time.Time
}

type FollowStatistic struct {
	Uid       int64
	Followers int64
	Following int64
	Followed  bool
}
