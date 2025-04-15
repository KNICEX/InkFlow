package repo

import (
	"context"
	"encoding/json"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/repo/dao"
	"github.com/samber/lo"
)

type NotificationRepo interface {
	CreateNotification(ctx context.Context, n domain.Notification) error
	DelByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
	FindByType(ctx context.Context, recipientId int64, types []domain.NotificationType, maxId int64, limit int) ([]domain.Notification, error)
	FindMergedLike(ctx context.Context, recipient int64, offset, limit int) ([]domain.MergedLikeNotification, error)
	DeleteMergedLike(ctx context.Context, recipient int64, subjectType domain.SubjectType, subjectId int64) error
	CountUnreadByType(ctx context.Context, recipientId int64, types []domain.NotificationType) (map[domain.NotificationType]int64, error)
	CountTotalUnread(ctx context.Context, recipientId int64) (int64, error)
	MarkAllRead(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
	Delete(ctx context.Context, recipientId int64, ids []int64) error
	DeleteByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
}

type NoCacheNotificationRepo struct {
	dao dao.NotificationDAO
}

func NewNoCacheNotificationRepo(dao dao.NotificationDAO) NotificationRepo {
	return &NoCacheNotificationRepo{
		dao: dao,
	}
}

func (repo *NoCacheNotificationRepo) CreateNotification(ctx context.Context, n domain.Notification) error {
	entity, err := repo.toEntity(n)
	if err != nil {
		return err
	}
	return repo.dao.Insert(ctx, entity)
}

func (repo *NoCacheNotificationRepo) DelByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return repo.dao.DeleteAll(ctx, recipientId, lo.Map(types, func(item domain.NotificationType, index int) string {
		return item.ToString()
	})...)
}

func (repo *NoCacheNotificationRepo) FindByType(ctx context.Context, recipientId int64, types []domain.NotificationType, maxId int64, limit int) ([]domain.Notification, error) {
	notifications, err := repo.dao.FindByType(ctx, recipientId, lo.Map(types, func(item domain.NotificationType, index int) string {
		return item.ToString()
	}), maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(notifications, func(item dao.Notification, index int) domain.Notification {
		domainNotification, _ := repo.toDomain(item)
		return domainNotification
	}), nil
}

func (repo *NoCacheNotificationRepo) FindMergedLike(ctx context.Context, recipient int64, offset, limit int) ([]domain.MergedLikeNotification, error) {
	ml, err := repo.dao.FindMergedLike(ctx, recipient, offset, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(ml, func(item dao.MergedLike, index int) domain.MergedLikeNotification {
		return domain.MergedLikeNotification{
			UserIds:     item.UserIds,
			Total:       item.Total,
			SubjectType: domain.SubjectTypeFromStr(item.SubjectType),
			SubjectId:   item.SubjectId,
			Read:        item.Read,
			UpdatedAt:   item.UpdatedAt,
		}
	}), nil
}

func (repo *NoCacheNotificationRepo) DeleteMergedLike(ctx context.Context, recipient int64, subjectType domain.SubjectType, subjectId int64) error {
	return repo.dao.DeleteMergedLike(ctx, recipient, subjectType.ToString(), subjectId)
}

func (repo *NoCacheNotificationRepo) CountUnreadByType(ctx context.Context, recipientId int64, types []domain.NotificationType) (map[domain.NotificationType]int64, error) {
	counts, err := repo.dao.CountUnreadByType(ctx, recipientId, lo.Map(types, func(item domain.NotificationType, index int) string {
		return item.ToString()
	}))
	if err != nil {
		return nil, err
	}

	return lo.MapEntries(counts, func(key string, value int64) (domain.NotificationType, int64) {
		return domain.NotificationTypeFromStr(key), value
	}), nil
}

func (repo *NoCacheNotificationRepo) CountTotalUnread(ctx context.Context, recipientId int64) (int64, error) {
	count, err := repo.dao.CountTotalUnread(ctx, recipientId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *NoCacheNotificationRepo) MarkAllRead(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return repo.dao.ReadAll(ctx, recipientId, lo.Map(types, func(item domain.NotificationType, index int) string {
		return item.ToString()
	})...)
}

func (repo *NoCacheNotificationRepo) Delete(ctx context.Context, recipientId int64, ids []int64) error {
	return repo.dao.Delete(ctx, ids, recipientId)
}

func (repo *NoCacheNotificationRepo) DeleteByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return repo.dao.DeleteAll(ctx, recipientId, lo.Map(types, func(item domain.NotificationType, index int) string {
		return item.ToString()
	})...)
}

func (repo *NoCacheNotificationRepo) toEntity(n domain.Notification) (dao.Notification, error) {
	content, err := json.Marshal(n.Content)
	if err != nil {
		return dao.Notification{}, err
	}
	return dao.Notification{
		Id:               n.Id,
		RecipientId:      n.RecipientId,
		SenderId:         n.SenderId,
		NotificationType: n.NotificationType.ToString(),
		SubjectType:      n.SubjectType.ToString(),
		SubjectId:        n.SubjectId,
		Content:          string(content),
		Read:             n.Read,
		CreatedAt:        n.CreatedAt,
	}, nil
}

func (repo *NoCacheNotificationRepo) toDomain(n dao.Notification) (domain.Notification, error) {
	var content any
	switch n.SubjectType {
	case domain.SubjectTypeComment.ToString():
		content = domain.ReplyContent{}
		if err := json.Unmarshal([]byte(n.Content), &content); err != nil {
			return domain.Notification{}, err
		}
	// TODO 更多类型
	default:
		content = n.Content
	}

	return domain.Notification{
		Id:               n.Id,
		RecipientId:      n.RecipientId,
		SenderId:         n.SenderId,
		SubjectType:      domain.SubjectTypeFromStr(n.SubjectType),
		SubjectId:        n.SubjectId,
		NotificationType: domain.NotificationTypeFromStr(n.NotificationType),
		Content:          content,
		Read:             n.Read,
		CreatedAt:        n.CreatedAt,
	}, nil
}
