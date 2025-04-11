package event

import "time"

type UserCreateEvent struct {
	UserId    int64     `json:"userId"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	AboutMe   string    `json:"aboutMe"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserUpdateEvent struct {
	UserId    int64     `json:"userId"`
	Username  string    `json:"username"`
	Account   string    `json:"account"`
	Avatar    string    `json:"avatar"`
	AboutMe   string    `json:"aboutMe"`
	Banner    string    `json:"banner"`
	Birthday  time.Time `json:"birthday"`
	Links     []string  `json:"links"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
