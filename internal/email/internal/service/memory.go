package service

import (
	"context"
	"fmt"
)

type MemoryService struct {
}

func NewMemoryService() Service {
	return &MemoryService{}
}

func (m MemoryService) SendString(ctx context.Context, email, title string, body string) error {
	fmt.Println("email:", email)
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	return nil
}

func (m MemoryService) SendHTML(ctx context.Context, email, title string, body string) error {
	fmt.Println("email:", email)
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	return nil
}

func (m MemoryService) Ping(ctx context.Context) error {
	return nil
}
