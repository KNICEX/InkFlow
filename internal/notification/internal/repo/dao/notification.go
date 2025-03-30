package dao

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// Notification
// 取消点赞后，又重新点赞
// 删除评论后，又重新评论...
type Notification struct {
	Id               int64
	RecipientId      int64
	SenderId         int64
	NotificationType string
	SubjectType      string
	SubjectId        int64
	Content          string
	CreatedAt        time.Time `gorm:"index:created_read"`
	Read             bool      `gorm:"index:created_read"`
}

type MergedLike struct {
	UserIds     []int64
	Total       int64
	SubjectType string
	SubjectId   int64
	Read        bool
	UpdatedAt   time.Time
}

type NotificationDAO interface {
	Insert(ctx context.Context, no Notification) error
	Delete(ctx context.Context, ids []int64, uid int64) error
	DeleteAll(ctx context.Context, uid int64, notificationType ...string) error
	FindByType(ctx context.Context, uid int64, notificationType []string, maxId int64, limit int) ([]Notification, error)
	FindLikeMerge(ctx context.Context, uid int64, offset, limit int) ([]MergedLike, error)
	ReadAll(ctx context.Context, userId int64, notificationType ...string) error
	CountTotalUnread(ctx context.Context, userId int64) (int64, error)
	CountUnreadByType(ctx context.Context, userId int64, types []string) (map[string]int64, error)
}

type GormNotificationDAO struct {
	node snowflakex.Node
	db   *gorm.DB
}

func NewGormNotificationDAO(db *gorm.DB, node snowflakex.Node) NotificationDAO {
	return &GormNotificationDAO{
		node: node,
		db:   db,
	}
}

func (dao *GormNotificationDAO) Insert(ctx context.Context, no Notification) error {
	no.Id = dao.node.NextID()
	no.CreatedAt = time.Now()
	return dao.db.WithContext(ctx).Create(&no).Error
}

func (dao *GormNotificationDAO) FindByType(ctx context.Context, uid int64, notificationType []string, maxId int64, limit int) ([]Notification, error) {
	var notifications []Notification
	err := dao.db.WithContext(ctx).Where("recipient_id = ? AND notification_type IN ? AND id < ?", uid, notificationType, maxId).
		Order("id desc").Limit(limit).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (dao *GormNotificationDAO) FindLikeMerge(ctx context.Context, uid int64, offset, limit int) ([]MergedLike, error) {
	// 按subject_type和subject_id分组，查出前n个用户id和总数量
	type subject struct {
		SubjectType string
		SubjectId   int64
		Total       int64
		UpdatedAt   time.Time
	}
	var subjects []subject

	err := dao.db.WithContext(ctx).Model(&Notification{}).
		Select("subject_type, subject_id, count(*) as total, MAX(created_at) as updated_at").
		Where("recipient_id = ? AND notification_type = ?", uid, "like").
		Group("subject_type, subject_id").
		Order("MAX(created_at) desc").
		Offset(offset).Limit(limit).Find(&subjects).Error
	if err != nil {
		return nil, err
	}
	if len(subjects) == 0 {
		return nil, nil
	}

	subjectInSql := strings.Builder{}
	subjectInSql.WriteString("(")
	for i, item := range subjects {
		subjectInSql.WriteString(fmt.Sprintf("('%s', ", item.SubjectType))
		subjectInSql.WriteString(strconv.FormatInt(item.SubjectId, 10))
		subjectInSql.WriteString(")")
		if i != len(subjects)-1 {
			subjectInSql.WriteString(", ")
		}
	}
	subjectInSql.WriteString(")")

	var recentUsers []Notification
	// 查出每个subject对应的前3个用户id
	subQuery := dao.db.Model(&Notification{}).
		Select("sender_id, subject_type, subject_id, read, created_at, "+
			"RANK() OVER (PARTITION BY subject_type, subject_id ORDER BY created_at DESC)").
		Where("recipient_id = ? AND notification_type = ?", uid, "like").
		Where("(subject_type, subject_id) IN " + subjectInSql.String())

	err = dao.db.WithContext(ctx).Table("(?) as t", subQuery).
		Select("sender_id, subject_type, subject_id, read").
		Where("rank <= 3").
		Order("created_at desc").
		Find(&recentUsers).Error
	if err != nil {
		return nil, err
	}
	subjectUserMap := make(map[string][]Notification)
	for _, item := range recentUsers {
		key := fmt.Sprintf("%s_%d", item.SubjectType, item.SubjectId)
		subjectUserMap[key] = append(subjectUserMap[key], item)
	}

	res := make([]MergedLike, 0, len(subjects))
	for _, item := range subjects {
		ml := MergedLike{
			SubjectType: item.SubjectType,
			SubjectId:   item.SubjectId,
			Total:       item.Total,
			Read:        true,
			UpdatedAt:   item.UpdatedAt,
		}
		key := fmt.Sprintf("%s_%d", item.SubjectType, item.SubjectId)
		if users, ok := subjectUserMap[key]; ok {
			ml.UserIds = make([]int64, 0, len(users))
			for _, user := range users {
				if !user.Read {
					ml.Read = false
				}
				ml.UserIds = append(ml.UserIds, user.SenderId)
			}
		}
		res = append(res, ml)
	}
	return res, err
}

func (dao *GormNotificationDAO) ReadAll(ctx context.Context, userId int64, notificationType ...string) error {
	if len(notificationType) == 0 {
		return dao.db.WithContext(ctx).Model(&Notification{}).Where("recipient_id = ?", userId).Update("read", true).Error
	}
	return dao.db.WithContext(ctx).Model(&Notification{}).
		Where("recipient_id = ? AND notification_type IN ?", userId, notificationType).
		Update("read", true).Error
}

func (dao *GormNotificationDAO) Delete(ctx context.Context, ids []int64, uid int64) error {
	return dao.db.WithContext(ctx).Where("recipient_id = ? AND id IN ?", uid, ids).Delete(&Notification{}).Error
}

func (dao *GormNotificationDAO) DeleteAll(ctx context.Context, uid int64, notificationType ...string) error {
	if len(notificationType) == 0 {
		return dao.db.WithContext(ctx).Where("recipient_id = ?", uid).Delete(&Notification{}).Error
	}
	return dao.db.WithContext(ctx).Where("recipient_id = ? AND notification_type IN ?", uid, notificationType).Delete(&Notification{}).Error
}

func (dao *GormNotificationDAO) CountTotalUnread(ctx context.Context, userId int64) (int64, error) {
	var cnt int64
	err := dao.db.WithContext(ctx).Model(&Notification{}).Where("recipient_id = ? AND read = false", userId).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (dao *GormNotificationDAO) CountUnreadByType(ctx context.Context, userId int64, types []string) (map[string]int64, error) {
	type UnreadCount struct {
		NotificationType string
		Cnt              int64
	}
	var notifications []UnreadCount
	err := dao.db.WithContext(ctx).Model(&Notification{}).
		Select("notification_type, count(*) as cnt").
		Where("recipient_id = ? AND read = false AND notification_type IN ?", userId, types).
		Group("notification_type").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64)
	for _, item := range notifications {
		res[item.NotificationType] = item.Cnt
	}
	return res, nil
}
