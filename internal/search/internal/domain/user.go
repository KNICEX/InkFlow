package domain

import "time"

type User struct {
	Id          int64
	Account     string
	Username    string
	FollowerCnt int64
	CreatedAt   time.Time
}

type UserOrderField string

const (
	UserOrderTypeDefault   UserOrderField = "default"
	UserOrderTypeFollower  UserOrderField = "follower"
	UserOrderTypeCreatedAt UserOrderField = "created_at"
)

type UserOrder struct {
	Field UserOrderField
	Desc  bool
}
