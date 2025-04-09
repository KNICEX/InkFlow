package service

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/code/internal/repo"
	"github.com/KNICEX/InkFlow/internal/email"
	"html/template"
	"math/rand"
	"time"
)

var (
	ErrCodeSendTooMany = repo.ErrCodeSendTooMany
	ErrCodeVerifyLimit = repo.ErrCodeVerifyLimit
)

//go:embed template.html
var defaultTemplate string

type Service interface {
	Send(ctx context.Context, biz, recipient string) error
	Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error)
}

type CachedEmailCodeService struct {
	repo     repo.CodeRepo
	emailSvc email.Service

	title    string
	template *template.Template // 密码将使用{code}替换

	effectiveTime  time.Duration
	resendInterval time.Duration
	maxRetry       int
}

func NewCachedEmailCodeService(repo repo.CodeRepo, emailSvc email.Service, opts ...CodeServiceOption) Service {
	temp, err := template.New("code").Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}
	svc := &CachedEmailCodeService{
		repo:     repo,
		emailSvc: emailSvc,

		title:    "InkFlow",
		template: temp,

		effectiveTime:  time.Minute * 5,
		resendInterval: time.Second * 10,
		maxRetry:       3,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

type CodeServiceOption func(option *CachedEmailCodeService)

func WithEffectiveTime(effectiveTime time.Duration) CodeServiceOption {
	return func(option *CachedEmailCodeService) {
		option.effectiveTime = effectiveTime
	}
}
func WithResendInterval(resendInterval time.Duration) CodeServiceOption {
	return func(option *CachedEmailCodeService) {
		option.resendInterval = resendInterval
	}
}

func WithMaxRetry(maxRetry int) CodeServiceOption {
	return func(option *CachedEmailCodeService) {
		option.maxRetry = maxRetry
	}
}

func WithTemplate(title string, temp *template.Template) CodeServiceOption {
	return func(option *CachedEmailCodeService) {
		option.title = title
		option.template = temp
	}
}

func (c *CachedEmailCodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}

func (c *CachedEmailCodeService) Send(ctx context.Context, biz, recipient string) error {
	code := c.generateCode()
	if err := c.repo.Store(ctx, biz, recipient, code, c.effectiveTime, c.resendInterval, c.maxRetry); err != nil {
		return err
	}
	var tempBuf bytes.Buffer
	if err := c.template.Execute(&tempBuf, map[string]string{
		"code": code,
	}); err != nil {
		return err
	}
	return c.emailSvc.SendHTML(ctx, recipient, c.title, tempBuf.String())
}

func (c *CachedEmailCodeService) Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error) {
	return c.repo.Verify(ctx, biz, recipient, inputCode)
}
