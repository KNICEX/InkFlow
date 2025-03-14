package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/uuidx"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAccountOrPwd = errors.New("账号/邮箱/手机号或密码错误")
	ErrUserDuplicate       = repo.ErrUserDuplicate
	ErrUserNotFound        = repo.ErrUserNotFound
)

type UserService interface {
	LoginEmailPwd(ctx context.Context, email, password string) (domain.User, error)
	LoginPhonePwd(ctx context.Context, phone, password string) (domain.User, error)
	LoginAccountPwd(ctx context.Context, account string, password string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error
	UpdateAccountName(ctx context.Context, uid int64, accountName string) error

	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByIds(ctx context.Context, uids []int64) (map[int64]domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByAccount(ctx context.Context, accountName string) (domain.User, error)
	FindByGithubId(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)

	ResetPwd(ctx context.Context, uid int64, newPwd string) error
	ChangePwd(ctx context.Context, uid int64, oldPwd, newPwd string) error
}

type userService struct {
	repo   repo.UserRepo
	l      logx.Logger
	tracer trace.Tracer
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
	return fmt.Sprintf("inker_%s", uuidx.NewShort())
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
	user, err := svc.repo.FindByAccount(ctx, accountName)
	if err != nil {
		return domain.User{}, err
	}
	if !svc.checkPwd(user.Password, password) {
		return domain.User{}, ErrInvalidAccountOrPwd
	}
	return user, nil
}
func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (svc *userService) FindByIds(ctx context.Context, uids []int64) (map[int64]domain.User, error) {
	users, err := svc.repo.FindByIds(ctx, uids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.User)
	for _, user := range users {
		res[user.Id] = user
	}
	return res, nil
}

func (svc *userService) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return svc.repo.FindByEmail(ctx, email)
}

func (svc *userService) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindByGithubId(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindByGithubId(ctx, id)
}

func (svc *userService) FindByAccount(ctx context.Context, accountName string) (domain.User, error) {
	return svc.repo.FindByAccount(ctx, accountName)
}

func (svc *userService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	// 理论上走到这里都是通过两步注册流程
	encryptedPwd, err := svc.encryptPwd(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = encryptedPwd
	if user.Id, err = svc.repo.Create(ctx, user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// 不使用后端自动创建的方法，让用户自己设置账户id，且后续不再允许变更
//func (svc *userService) FindOrCreateByGithub(ctx context.Context, i domain.GithubInfo) (domain.User, error) {
//	user, err := svc.repo.FindByGithubId(ctx, i.Id)
//	if err == nil {
//		return user, nil
//	}
//	if !errors.Is(err, repo.ErrUserNotFound) {
//		return domain.User{}, err
//	}
//	user = domain.User{
//		Account:    svc.newDefaultAccountName(),
//		GithubInfo: i,
//	}
//	if user.Id, err = svc.repo.Create(ctx, user); err != nil {
//		return domain.User{}, err
//	}
//	return user, nil
//}

func (svc *userService) UpdateAccountName(ctx context.Context, uid int64, accountName string) error {
	// TODO 这个接口应该不会开放给用户
	return svc.repo.UpdateNonZeroFields(ctx, domain.User{
		Id:      uid,
		Account: accountName,
	})
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
