package user

import (
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/oauth2"
)

var (
	ErrInvalidAccountOrPwd = service.ErrInvalidAccountOrPwd
	ErrUserDuplicate       = service.ErrUserDuplicate
)

type Service = service.UserService
type OAuth2Service[T any] = oauth2.Service[T]

type User = domain.User
type GithubInfo = domain.GithubInfo
