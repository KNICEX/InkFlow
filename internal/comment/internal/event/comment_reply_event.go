package event

import "time"

// CommentReplyEvent represents an event triggered when a user replies to a comment
type CommentReplyEvent struct {
	Id             int64     `json:"id"`
	UserId         int64     `json:"user_id"`
	CommentId      int64     `json:"comment_id"`
	ReplyUserId    int64     `json:"reply_user_id"`
	ReplyCommentId int64     `json:"reply_comment_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
