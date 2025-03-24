package domain

import "time"

type Notification struct {
	Id               int64
	RecipientId      int64
	SenderId         int64
	SubjectType      SubjectType
	SubjectId        int64
	NotificationType NotificationType
	Content          any
	Read             bool
	CreatedAt        time.Time
}

type NotificationType string

const (
	NotificationTypeLike      NotificationType = "like"
	NotificationTypeReply     NotificationType = "reply"
	NotificationTypeSubscribe NotificationType = "subscribe"
	NotificationTypeMention   NotificationType = "mention"
	NotificationTypeFollow    NotificationType = "follow"
	NotificationTypeSystem    NotificationType = "system"
)

func (t NotificationType) ToString() string {
	return string(t)
}

func NotificationTypeFromStr(s string) NotificationType {
	return NotificationType(s)
}

type SubjectType string

const (
	SubjectTypeComment SubjectType = "comment"
	SubjectTypeInk     SubjectType = "ink"
	SubjectTypeUser    SubjectType = "user"
	SubjectTypeFeed    SubjectType = "feed"
	SubjectTypeSystem  SubjectType = "system"
)

func (t SubjectType) ToString() string {
	return string(t)
}

func SubjectTypeFromStr(s string) SubjectType {
	return SubjectType(s)
}

type ReplyContent struct {
	RootId        int64
	ParentId      int64
	ReplyContent  string // 回复内容
	TargetId      int64  // 被回复的评论id, 如果是根评论则为0
	TargetContent string // 被回复的内容, 如果是根评论则为空
	//RootContent   string // 根评论内容
}

type LikeContent struct {
}
