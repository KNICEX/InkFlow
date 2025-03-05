package service

import "context"

type Service interface {
	SendString(ctx context.Context, email, title string, body string) error
	SendHTML(ctx context.Context, email, title string, body string) error
	Ping(ctx context.Context) error
}
