package gorse

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/service"
	client "github.com/gorse-io/gorse-go"
	"strconv"
	"time"
)

type SyncService struct {
	cli *client.GorseClient
}

func NewSyncService(cli *client.GorseClient) service.SyncService {
	return &SyncService{cli: cli}
}

func (s SyncService) InputUser(ctx context.Context, user domain.User) error {
	_, err := s.cli.InsertUser(ctx, client.User{
		UserId: strconv.FormatInt(user.Id, 10),
	})
	return err
}

func (s SyncService) InputInk(ctx context.Context, ink domain.Ink) error {
	_, err := s.cli.InsertItem(ctx, client.Item{
		ItemId:    strconv.FormatInt(ink.Id, 10),
		Labels:    ink.Tags,
		Timestamp: ink.CreatedAt.Format(time.RFC3339),
		Comment:   ink.Title,
	})
	return err
}

func (s SyncService) InputFeedback(ctx context.Context, feedback domain.Feedback) error {
	_, err := s.cli.InsertFeedback(ctx, []client.Feedback{
		{
			UserId:       strconv.FormatInt(feedback.UserId, 10),
			ItemId:       strconv.FormatInt(feedback.InkId, 10),
			FeedbackType: feedback.FeedbackType.ToString(),
			Timestamp:    feedback.CreatedAt.Format(time.RFC3339),
		},
	})
	return err
}

func (s SyncService) InputRelation(ctx context.Context, relation domain.Relation) error {
	return nil
}

func (s SyncService) DeleteUser(ctx context.Context, userId int64) error {
	_, err := s.cli.DeleteUser(ctx, strconv.FormatInt(userId, 10))
	return err
}

func (s SyncService) HiddenInk(ctx context.Context, inkId int64) error {
	hidden := true
	_, err := s.cli.UpdateItem(ctx, strconv.FormatInt(inkId, 10), client.ItemPatch{
		IsHidden: &hidden,
	})
	return err
}
