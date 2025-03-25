package domain

import "time"

type User struct {
	Id        int64
	Avatar    string
	Account   string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
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
