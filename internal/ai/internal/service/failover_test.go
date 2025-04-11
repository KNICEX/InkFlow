package service

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
	"testing"
)

type mockLLMService struct {
	name string
}

func (m mockLLMService) AskOnce(ctx context.Context, question string) (domain.Resp, error) {
	fmt.Println("mockLLMService", m.name, "question:", question)
	return domain.Resp{}, fmt.Errorf("mock err %s", m.name)
}

func (m mockLLMService) BeginChat(ctx context.Context) (LLMSession, error) {
	//TODO implement me
	panic("implement me")
}

func TestFailoverLLMService_AskOnce(t *testing.T) {
	m1 := mockLLMService{name: "m1"}
	m2 := mockLLMService{name: "m2"}
	m3 := mockLLMService{name: "m3"}
	m4 := mockLLMService{name: "m4"}

	svc := NewFailoverService([]LLMService{m1, m2, m3, m4})

	_, _ = svc.AskOnce(context.Background(), "test")

}
