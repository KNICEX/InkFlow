package code

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/repo"
)

var (
	ErrCodeSendTooMany = repo.ErrCodeSendTooMany
	ErrCodeVerifyLimit = repo.ErrCodeVerifyLimit
)

type Service interface {
	Send(ctx context.Context, biz, recipient string) error
	Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error)
}
