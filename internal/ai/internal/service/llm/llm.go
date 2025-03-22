package llm

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
)

type Session interface {
	Ask(ctx context.Context, question string) (domain.Resp, error)
	Close() error
}

type Service interface {
	AskOnce(ctx context.Context, question string) (domain.Resp, error)
	BeginChat(ctx context.Context) (Session, error)
}
