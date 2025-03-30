package notification

import (
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
)

type Service = service.NotificationService

type Notification = domain.Notification
type MergedLike = domain.MergedLikeNotification

type Type = domain.NotificationType
type SubjectType = domain.SubjectType

const (
	TypeReply     = domain.NotificationTypeReply
	TypeLike      = domain.NotificationTypeLike
	TypeFollow    = domain.NotificationTypeFollow
	TypeMention   = domain.NotificationTypeMention
	TypeSubscribe = domain.NotificationTypeSubscribe
	TypeSystem    = domain.NotificationTypeSystem

	SubjectTypeInk     = domain.SubjectTypeInk
	SubjectTypeComment = domain.SubjectTypeComment
)

type ReplyContent = domain.ReplyContent
