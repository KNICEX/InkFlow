package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/comment"
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

func (u *userAggregate) GetUser(ctx context.Context, uid int64, viewUid int64) (UserVO, error) {
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
	inkSvc        ink.Service
	userAggregate *userAggregate
	intrAggregate *interactiveAggregate
}

func newInkAggregate(inkSvc ink.Service, userSvc user.Service, followSvc relation.FollowService,
	intrSvc interactive.Service, commentSvc comment.Service) *inkAggregate {
	return &inkAggregate{
		inkSvc:        inkSvc,
		userAggregate: newUserAggregate(userSvc, followSvc),
		intrAggregate: newInteractiveAggregate(intrSvc, commentSvc),
	}
}

func (i *inkAggregate) GetInk(ctx context.Context, id int64, viewUid int64) (InkVO, error) {
	var author UserVO
	var intr InteractiveVO
	inkInfo, err := i.inkSvc.FindById(ctx, id)
	if err != nil {
		return InkVO{}, err
	}
	if inkInfo.Id == 0 {
		return InkVO{}, nil
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		author, er = i.userAggregate.GetUser(ctx, inkInfo.Author.Id, viewUid)
		return er
	})
	eg.Go(func() error {
		var er error
		intr, er = i.intrAggregate.GetInteractive(ctx, bizInk, inkInfo.Id, viewUid)
		return er
	})
	if err = eg.Wait(); err != nil {
		return InkVO{}, err
	}

	vo := inkToVO(inkInfo)
	vo.Author = author
	vo.Interactive = intr
	return vo, nil
}

func (i *inkAggregate) GetInkList(ctx context.Context, ids []int64, viewUid int64) (map[int64]InkVO, error) {
	var authors map[int64]UserVO
	var intrs map[int64]InteractiveVO
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
		authors, er = i.userAggregate.GetUserList(ctx, authorIds, viewUid)
		return er
	})
	eg.Go(func() error {
		var er error
		intrs, er = i.intrAggregate.GetInteractiveList(ctx, bizInk, ids, viewUid)
		return er
	})
	if err = eg.Wait(); err != nil {
		return nil, err
	}

	vos := make(map[int64]InkVO, len(inkMap))
	for _, inkInfo := range inkMap {
		vo := inkToVO(inkInfo)
		vo.Author = authors[inkInfo.Author.Id]
		vo.Interactive = intrs[inkInfo.Id]
		vos[inkInfo.Id] = vo
	}
	return vos, nil
}

type interactiveAggregate struct {
	intrSvc    interactive.Service
	commentSvc comment.Service
}

func newInteractiveAggregate(intrSvc interactive.Service, commentSvc comment.Service) *interactiveAggregate {
	return &interactiveAggregate{
		intrSvc:    intrSvc,
		commentSvc: commentSvc,
	}
}

func (i *interactiveAggregate) GetInteractive(ctx context.Context, biz string, id int64, uid int64) (InteractiveVO, error) {
	var intr interactive.Interactive
	var commentCounts map[int64]int64
	eg := errgroup.Group{}

	eg.Go(func() error {
		var err error
		intr, err = i.intrSvc.Get(ctx, biz, id, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		commentCounts, err = i.commentSvc.FindBizReplyCount(ctx, biz, []int64{id})
		return err
	})
	if err := eg.Wait(); err != nil {
		return InteractiveVO{}, err
	}
	vo := intrToVo(intr)
	vo.CommentCnt = commentCounts[id]
	return vo, nil
}

func (i *interactiveAggregate) GetInteractiveList(ctx context.Context, biz string, ids []int64, uid int64) (map[int64]InteractiveVO, error) {
	var intrs map[int64]interactive.Interactive
	var commentCounts map[int64]int64
	eg := errgroup.Group{}

	eg.Go(func() error {
		var err error
		intrs, err = i.intrSvc.GetMulti(ctx, biz, ids, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		commentCounts, err = i.commentSvc.FindBizReplyCount(ctx, biz, ids)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	vos := make(map[int64]InteractiveVO, len(intrs))
	for _, intr := range intrs {
		vo := intrToVo(intr)
		vo.CommentCnt = commentCounts[intr.BizId]
		vos[intr.BizId] = vo
	}
	return vos, nil
}
