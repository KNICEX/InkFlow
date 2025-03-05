package event

import "time"

type CommentLikeEvt struct {
	Id         int64     `json:"id"`
	CommentId  int64     `json:"comment_id"`
	UserId     int64     `json:"user_id"`
	LikeUserId int64     `json:"like_user_id"`
	CreatedAt  time.Time `json:"created_at"`
}
