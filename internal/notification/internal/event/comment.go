package event

type CommentInteractive struct {
	CommentId     int64
	CommentatorId int64
	Uid           int64
	Biz           string
	BizId         int64
	Type          string
	Ext           map[string]string
}

const (
	CommentLikeType  = "like"
	CommentReplyType = "reply"
)
