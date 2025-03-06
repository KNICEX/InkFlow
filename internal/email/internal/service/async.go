package service

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

// AsyncService TODO 考虑用kafka
type AsyncService struct {
	svc Service
	l   logx.Logger
}

func NewAsyncService(svc Service, l logx.Logger) Service {
	return &AsyncService{
		svc: svc,
		l:   l,
	}
}

func (a *AsyncService) SendString(ctx context.Context, email, title string, body string) error {
	go func() {
		if err := a.svc.SendString(ctx, email, title, body); err != nil {
			a.l.Error("failed to async send string email",
				logx.Error(err),
				logx.String("email", email))
		}
	}()
	return nil
}

func (a *AsyncService) SendHTML(ctx context.Context, email, title string, body string) error {
	go func() {
		if err := a.svc.SendHTML(ctx, email, title, body); err != nil {
			a.l.Error("failed to async send html email",
				logx.Error(err),
				logx.String("email", email))
		}
	}()
	return nil
}

func (a *AsyncService) Ping(ctx context.Context) error {
	return a.svc.Ping(ctx)
}
