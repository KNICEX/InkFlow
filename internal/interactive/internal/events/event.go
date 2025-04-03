package events

import "time"

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

type ReplyEvent struct {
	CommentId     int64     `json:"commentId"`
	RootId        int64     `json:"rootId"`
	ParentId      int64     `json:"parentId"`
	Biz           string    `json:"biz"`
	BizId         int64     `json:"bizId"`
	CommentatorId int64     `json:"commentatorId"`
	Payload       Payload   `json:"payload"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Payload struct {
	Content string   `json:"content"`
	Images  []string `json:"images"`
}
