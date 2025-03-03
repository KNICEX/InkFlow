package service

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/uuidx"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAccountOrPwd = errors.New("账号/邮箱/手机号或密码错误")
	ErrUserDuplicate       = repo.ErrUserDuplicate
)

type UserService interface {
	LoginEmailPwd(ctx context.Context, email, password string) (domain.User, error)
	LoginPhonePwd(ctx context.Context, phone, password string) (domain.User, error)
	LoginAccountPwd(ctx context.Context, accountName string, password string) (domain.User, error)
	Profile(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error
	UpdateAccountName(ctx context.Context, uid int64, accountName string) error
	FindOrCreateByGithub(ctx context.Context, i domain.GithubInfo) (domain.User, error)

	ResetPwd(ctx context.Context, uid int64, newPwd string) error
	ChangePwd(ctx context.Context, uid int64, oldPwd, newPwd string) error
	LoginAccountNamePwd(ctx *gin.Context, accountName string, password string) (domain.User, error)
}

type userService struct {
	repo repo.UserRepo
	l    logx.Logger
}

func NewUserService(repo repo.UserRepo, l logx.Logger) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

func (svc *userService) checkPwd(encrypted, plain string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plain)); err != nil {
		return false
	}
	return true
}

func (svc *userService) encryptPwd(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (svc *userService) newDefaultAccountName() string {
	return uuidx.NewShortN(12)
}
func (svc *userService) LoginEmailPwd(ctx context.Context, email, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if !svc.checkPwd(user.Password, password) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	return user, nil
}

func (svc *userService) LoginPhonePwd(ctx context.Context, phone, password string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	if !svc.checkPwd(user.Password, password) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	return user, nil
}
func (svc *userService) LoginAccountPwd(ctx context.Context, accountName string, password string) (domain.User, error) {
	user, err := svc.repo.FindByAccountName(ctx, accountName)
	if err != nil {
		return domain.User{}, err
	}
	if !svc.checkPwd(user.Password, password) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	return user, nil
}
func (svc *userService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}
func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}
func (svc *userService) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			user = domain.User{
				Account: svc.newDefaultAccountName(),
				Phone:   phone,
			}
			if err := svc.repo.Create(ctx, user); err != nil {
				return domain.User{}, err
			}
		} else {
			return domain.User{}, err
		}
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			user = domain.User{
				Account: svc.newDefaultAccountName(),
				Email:   email,
			}
			if err := svc.repo.Create(ctx, user); err != nil {
				return domain.User{}, err
			}
		} else {
			return domain.User{}, err
		}
	}
	return svc.repo.FindById(ctx, user.Id)
}

func (svc *userService) UpdateAccountName(ctx context.Context, uid int64, accountName string) error {
	// TODO 这个接口也许要做频率限制，不应当频繁修改
	return svc.repo.UpdateNonZeroFields(ctx, domain.User{
		Id:      uid,
		Account: accountName,
	})
}

func (svc *userService) FindOrCreateByGithub(ctx context.Context, i domain.GithubInfo) (domain.User, error) {
	user, err := svc.repo.FindByGithubId(ctx, i.Id)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			user = domain.User{
				Account:    svc.newDefaultAccountName(),
				GithubInfo: i,
			}
			if err := svc.repo.Create(ctx, user); err != nil {
				return domain.User{}, err
			}
		} else {
			return domain.User{}, err
		}
	}
	return svc.repo.FindByGithubId(ctx, i.Id)
}

func (svc *userService) ResetPwd(ctx context.Context, uid int64, newPwd string) error {
	encryptedPwd, err := svc.encryptPwd(newPwd)
	if err != nil {
		return err
	}
	return svc.repo.UpdateNonZeroFields(ctx, domain.User{
		Id:       uid,
		Password: encryptedPwd,
	})
}

func (svc *userService) ChangePwd(ctx context.Context, uid int64, oldPwd, newPwd string) error {
	u, err := svc.repo.FindById(ctx, uid)
	if err != nil {
		return err
	}
	if !svc.checkPwd(u.Password, oldPwd) {
		return ErrInvalidAccountOrPwd
	}
	newPwd, err = svc.encryptPwd(newPwd)
	if err != nil {
		return err
	}
	return svc.repo.UpdateNonZeroFields(ctx, domain.User{
		Id:       uid,
		Password: newPwd,
	})
}

func (svc *userService) LoginAccountNamePwd(ctx *gin.Context, accountName string, password string) (domain.User, error) {
	u, err := svc.repo.FindByAccountName(ctx, accountName)
	if errors.Is(err, repo.ErrUserNotFound) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	if err != nil {
		return domain.User{}, err
	}
	if !svc.checkPwd(u.Password, password) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	return u, nil
}
