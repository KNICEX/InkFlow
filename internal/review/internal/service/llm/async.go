package llm

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
)

type AsyncService struct {
	llm ai.LLMService
}

func (a AsyncService) SubmitInk(ctx context.Context, ink domain.Ink) error {
	//TODO implement me
	panic("implement me")
}
