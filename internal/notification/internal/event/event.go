package event

import "time"

type FollowEvent struct {
	FollowerId int64
	FolloweeId int64
	CreatedAt  time.Time
}

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

type CommentLikeEvt struct {
	CommentId int64     `json:"commentId"`
	LikeUid   int64     `json:"likeUid"`
	CreatedAt time.Time `json:"createdAt"`
}

type InkLikeEvent struct {
	InkId     int64     `json:"inkId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
