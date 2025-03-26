package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
)

type LLMSession interface {
	Ask(ctx context.Context, question string) (domain.Resp, error)
	Close() error
}

type LLMService interface {
	AskOnce(ctx context.Context, question string) (domain.Resp, error)
	BeginChat(ctx context.Context) (LLMSession, error)
}
