package service

import "context"

type Resp struct {
	Content string
	Token   int
}

type Context struct {
	Resp Resp
	Err  error
	Ask  func(ctx context.Context, msg string) error
}

type LLMService interface {
	AskOnce(ctx context.Context, msg string) (Resp, error)
	// 连续对话
	BeginChat(ctx context.Context, msg string) <-chan Context
}
