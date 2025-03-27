package event

import "time"

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

type DeleteEvent struct {
	CommentId int64     `json:"commentId"`
	CreatedAt time.Time `json:"createdAt"`
}

type LikeEvent struct {
	CommentId int64     `json:"commentId"`
	LikeUid   int64     `json:"likeUid"`
	CreatedAt time.Time `json:"createdAt"`
}
