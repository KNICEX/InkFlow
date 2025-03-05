package domain

import "time"

type FollowRelation struct {
	Follower  int64
	Followee  int64
	CreatedAt time.Time
}

type FollowStatistic struct {
	Followers  int64
	Followings int64
}
