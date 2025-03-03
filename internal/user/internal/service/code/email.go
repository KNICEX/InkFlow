package code

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/email"
	"github.com/ecodeclub/ekit/bean/option"
	"math/rand"
	"time"
)

func NewEmailCodeService(codeRepo repo.CodeRepo, emailSvc email.Service, option option.Option[EmailCodeService]) Service {
	return &EmailCodeService{
		codeRepo: codeRepo,
		emailSvc: emailSvc,

		effectiveTime:  time.Minute * 5,
		resendInterval: time.Second * 10,
		maxRetry:       3,
	}
}

func WithEffectiveTime(effectiveTime time.Duration) option.Option[EmailCodeService] {
	return func(e *EmailCodeService) {
		e.effectiveTime = effectiveTime
	}
}

func WithResendInterval(resendInterval time.Duration) option.Option[EmailCodeService] {
	return func(e *EmailCodeService) {
		e.resendInterval = resendInterval
	}
}

func WithMaxRetry(maxRetry int) option.Option[EmailCodeService] {
	return func(e *EmailCodeService) {
		e.maxRetry = maxRetry
	}
}

type EmailCodeService struct {
	codeRepo repo.CodeRepo
	emailSvc email.Service

	effectiveTime  time.Duration
	resendInterval time.Duration
	maxRetry       int
}

func (e *EmailCodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

func (e *EmailCodeService) Send(ctx context.Context, biz, email string) error {
	code := e.generateCode()
	err := e.codeRepo.Store(ctx, biz, email, code, e.effectiveTime, e.resendInterval, e.maxRetry)
	if err != nil {
		return err
	}
	// TODO 自定义邮件模板
	return e.emailSvc.SendHTML(ctx, email, "验证码", fmt.Sprintf("您的验证码是: %s", code))
}

func (e *EmailCodeService) Verify(ctx context.Context, biz, email, code string) (bool, error) {
	return e.codeRepo.Verify(ctx, biz, email, code)
}
