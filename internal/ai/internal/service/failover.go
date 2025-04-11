package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
	"github.com/KNICEX/InkFlow/pkg/backoff"
	"sync/atomic"
)

type FailoverLLMService struct {
	svcs   []LLMService
	idx    *atomic.Int32
	policy backoff.Policy
}

type Option func(*FailoverLLMService)

func NewFailoverService(services []LLMService, policy ...backoff.Policy) LLMService {
	svc := &FailoverLLMService{
		svcs: services,
		idx:  new(atomic.Int32),
	}
	if len(policy) > 0 {
		svc.policy = policy[0]
	} else {
		svc.policy = backoff.DefaultPolicy
	}
	return svc
}

func (f *FailoverLLMService) AskOnce(ctx context.Context, question string) (domain.Resp, error) {
	svc := f.svcs[f.idx.Load()%int32(len(f.svcs))]
	f.idx.Add(1)

	var resp domain.Resp
	var err error

	copyPolicy := f.policy
	copyPolicy.OnRetry = func(i int, err error) {
		// 每次重试都换svc
		svc = f.svcs[f.idx.Load()%int32(len(f.svcs))]
		f.idx.Add(1)
		if f.policy.OnRetry != nil {
			f.policy.OnRetry(i, err)
		}
	}

	fn := backoff.Wrap(func() error {
		resp, err = svc.AskOnce(ctx, question)
		if err != nil {
			return err
		}
		return nil
	}, copyPolicy)

	err = fn()
	return resp, err
}

func (f *FailoverLLMService) BeginChat(ctx context.Context) (LLMSession, error) {
	i := int(f.idx.Load())
	f.idx.Add(1)
	var session LLMSession
	var err error

	for range len(f.svcs) {
		// 最多遍历一遍所有svc
		session, err = f.svcs[i].BeginChat(ctx)
		if err == nil {
			return session, nil
		}
		i++
	}
	return session, err
}
