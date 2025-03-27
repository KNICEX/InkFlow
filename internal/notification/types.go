package notification

import (
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
)

type Service = service.NotificationService

type Notification = domain.Notification
type Type = domain.NotificationType

type ReplyContent = domain.ReplyContent
