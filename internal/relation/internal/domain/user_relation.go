package domain

import "time"

type Follow struct {
	Id       int64
	Follower Follower
	Follows  []Followee

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Follower struct {
	Id   int64
	Name string
}

type Followee struct {
	Id   int64
	Name string
}

type FollowStatistic struct {
	Followers  int64
	Followings int64
}
