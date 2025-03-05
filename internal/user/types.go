package user

import (
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
)

var (
	ErrInvalidAccountOrPwd = service.ErrInvalidAccountOrPwd
	ErrUserDuplicate       = service.ErrUserDuplicate
)

type Service = service.UserService

type User = domain.User
