package event

import "time"

type ReplyEvent struct {
	CommentId     int64  `json:"commentId"`
	RootId        int64  `json:"rootId"`
	ParentId      int64  `json:"parentId"`
	Biz           string `json:"biz"`
	BizId         int64  `json:"bizId"`
	CommentatorId int64  `json:"commentatorId"`
	Payload       struct {
		Content string   `json:"content"`
		Images  []string `json:"images"`
	} `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserCreateEvent struct {
	UserId    int64     `json:"userId"`
	Avatar    string    `json:"avatar"`
	AboutMe   string    `json:"aboutMe"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
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
