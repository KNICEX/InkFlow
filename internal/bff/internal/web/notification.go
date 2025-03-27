package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

type NotificationHandler struct {
	svc        notification.Service
	userSvc    user.Service
	inkSvc     ink.Service
	commentSvc comment.Service
	l          logx.Logger
}
