package dao

import "time"

type Notification struct {
	Id        int64
	UserId    int64
	Type      NotificationType
	Content   string
	Read      bool
	CreatedAt time.Time
}

type NotificationType int

const (
	NotificationTypeUnknown        NotificationType = iota
	NotificationTypeInkLike                         // 作品被点赞
	NotificationTypeInkComment                      // 作品被评论
	NotificationTypeInkRef                          // 作品被引用
	NotificationTypeCommentLike                     // 评论被点赞
	NotificationTypeCommentReply                    // 评论被回复
	NotificationTypeFollow                          // 被关注
	NotificationTypeFollowerNewInk                  // 关注的人有新的作品
	NotificationTypeRef                             // 被@
	NotificationSystem                              // 系统通知
)
