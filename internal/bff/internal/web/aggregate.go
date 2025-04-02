package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type userAggregate struct {
	userSvc   user.Service
	followSvc relation.FollowService
}

func newUserAggregate(userSvc user.Service, followSvc relation.FollowService) *userAggregate {
	return &userAggregate{
		userSvc:   userSvc,
		followSvc: followSvc,
	}
}

func (u *userAggregate) GetUserDetail(ctx context.Context, uid int64, viewUid int64) (UserVO, error) {
	var userInfo user.User
	var followInfo relation.FollowStatistic

	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		userInfo, err = u.userSvc.FindById(ctx, uid)
		return err
	})

	eg.Go(func() error {
		var err error
		followInfo, err = u.followSvc.FindFollowStats(ctx, uid, viewUid)
		return err
	})
	if err := eg.Wait(); err != nil {
		return UserVO{}, err
	}
	vo := userToVO(userInfo)
	vo.Followed = followInfo.Followed
	vo.Following = followInfo.Following
	vo.Followers = followInfo.Followers
	return vo, nil
}

func (u *userAggregate) GetUserList(ctx context.Context, uids []int64, viewUid int64) (map[int64]UserVO, error) {
	var users map[int64]user.User
	var followInfos map[int64]relation.FollowStatistic
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		users, err = u.userSvc.FindByIds(ctx, uids)
		return err
	})
	eg.Go(func() error {
		var err error
		followInfos, err = u.followSvc.FindFollowStatsBatch(ctx, uids, viewUid)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	vos := make(map[int64]UserVO, len(users))
	for _, userInfo := range users {
		vo := userToVO(userInfo)
		if followInfo, ok := followInfos[userInfo.Id]; ok {
			vo.Followed = followInfo.Followed
			vo.Following = followInfo.Following
			vo.Followers = followInfo.Followers
		}
		vos[userInfo.Id] = vo
	}
	return vos, nil
}

type inkAggregate struct {
	inkSvc  ink.Service
	userSvc user.Service
	intrSvc interactive.Service
}

func newInkAggregate(inkSvc ink.Service, userSvc user.Service, intrSvc interactive.Service) *inkAggregate {
	return &inkAggregate{
		inkSvc:  inkSvc,
		userSvc: userSvc,
		intrSvc: intrSvc,
	}
}

func (i *inkAggregate) GetInkList(ctx context.Context, ids []int64, viewUid int64) (map[int64]InkVO, error) {
	var authors map[int64]user.User
	var intrs map[int64]interactive.Interactive
	inkMap, err := i.inkSvc.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(inkMap) == 0 {
		return nil, nil
	}

	authorIds := lo.MapToSlice(inkMap, func(key int64, value ink.Ink) int64 {
		return value.Author.Id
	})
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		authors, er = i.userSvc.FindByIds(ctx, authorIds)
		return er
	})
	eg.Go(func() error {
		var er error
		intrs, er = i.intrSvc.GetMulti(ctx, bizInk, ids, viewUid)
		return er
	})
	if err = eg.Wait(); err != nil {
		return nil, err
	}

	vos := make(map[int64]InkVO, len(inkMap))
	for _, inkInfo := range inkMap {
		vo := inkToVO(inkInfo)
		if author, ok := authors[inkInfo.Author.Id]; ok {
			vo.Author = userToVO(author)
		}
		if intr, ok := intrs[inkInfo.Id]; ok {
			vo.Interactive = intrToVo(intr)
		}
		vos[inkInfo.Id] = vo
	}
	return vos, nil
}
