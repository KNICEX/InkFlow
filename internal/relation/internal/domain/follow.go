package domain

import "time"

type FollowRelation struct {
	Follower  int64
	Followee  int64
	CreatedAt time.Time
}

type FollowStatistic struct {
	Followers int64
	Following int64
	Followed  bool
}

type FollowInfo struct {
	Uid      int64
	Followed bool
}
