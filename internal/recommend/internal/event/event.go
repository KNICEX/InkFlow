package event

import "time"

type UserCreateEvent struct {
	UserId    int64     `json:"userId"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
