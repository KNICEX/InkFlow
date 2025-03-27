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
	CommentId     int64
	SourceContent ReplyPayload
	TargetContent ReplyPayload
}

type ReplyPayload struct {
	Content string
	Images  []string
}

// TODO 好像暂时不需要冗余

type LikeContent struct {
}

type FollowContent struct {
}

type SubscribeContent struct {
}

type SystemContent struct {
}

type MentionContent struct {
}
