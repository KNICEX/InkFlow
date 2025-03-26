package llm

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"time"
)

type Service struct {
	llm ai.LLMService
}

func NewLLMService(llm ai.LLMService) service.Service {
	return &Service{
		llm: llm,
	}
}

func (s *Service) ReviewInk(ctx context.Context, ink domain.Ink) (domain.ReviewResult, error) {
	//session, err := s.llm.BeginChat(ctx)
	//if err != nil {
	//	return domain.ReviewResult{}, err
	//}
	//defer session.Close()
	//
	//// TODO ask in once or in a chat
	//resp, err := session.Ask(ctx, ink.Content)
	//if err != nil {
	//	return domain.ReviewResult{}, err
	//}
	//result, err := s.parseResp(resp.Content)
	//if err != nil {
	//	return domain.ReviewResult{}, err
	//}
	//return result, nil

	time.Sleep(time.Second * 20)
	return domain.ReviewResult{
		Passed: true,
	}, nil
}

func (s *Service) parseResp(resp string) (domain.ReviewResult, error) {
	// TODO parse response
	return domain.ReviewResult{}, nil
}
