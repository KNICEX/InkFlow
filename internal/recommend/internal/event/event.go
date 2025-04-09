package event

import "time"

type UserCreateEvent struct {
	UserId    int64     `json:"userId"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type InkViewEvent struct {
	InkId     int64     `json:"inkId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type InkLikeEvent struct {
	InkId     int64     `json:"inkId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type InkCancelLikeEvent struct {
	InkId     int64     `json:"inkId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
