package oauth2

import (
	"context"
	"errors"
)

var ErrCodeInvalid = errors.New("oauth2: code invalid")

type Service[T any] interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (T, error)
}
