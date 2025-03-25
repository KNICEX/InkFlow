package schedule

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo/dao"
	"github.com/samber/lo"
)

type DeleteHiddenSchedule struct {
	pullDAO dao.PullFeedDAO
	pushDAO dao.PushFeedDAO
}

func (s *DeleteHiddenSchedule) DeleteHiddenPull(ctx context.Context) error {
	feeds, err := s.pullDAO.FindHidden(ctx, 300)
	if err != nil {
		return err
	}
	return s.pullDAO.BatchDelete(ctx, lo.Map(feeds, func(item dao.PullFeed, index int) int64 {
		return item.Id
	}))
}

func (s *DeleteHiddenSchedule) DeleteHiddenPush(ctx context.Context) error {
	feeds, err := s.pushDAO.FindHidden(ctx, 800)
	if err != nil {
		return err
	}
	return s.pushDAO.BatchDelete(ctx, lo.Map(feeds, func(item dao.PushFeed, index int) int64 {
		return item.Id
	}))
}
