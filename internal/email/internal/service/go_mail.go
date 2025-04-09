package service

import (
	"context"
	"fmt"
	"github.com/wneessen/go-mail"
)

type MailService struct {
	host     string
	port     int
	authType mail.SMTPAuthType
	password string
	username string
	from     string

	client *mail.Client
}

type MailServiceOption func(*MailService)

func NewMailService(host string, port int, username, password, fromName string) (Service, error) {
	svc := &MailService{
		host:     host,
		port:     port,
		authType: mail.SMTPAuthPlain,
		password: password,
		username: username, // 发件邮箱
		from:     fromName, // 发件人(自定义昵称)
	}

	c, err := mail.NewClient(svc.host,
		mail.WithPort(svc.port),
		mail.WithSSLPort(true),
		mail.WithSMTPAuth(svc.authType),
		mail.WithUsername(svc.username),
		mail.WithPassword(svc.password),
	)
	if err != nil {
		return nil, err
	}
	svc.client = c

	return svc, nil
}

func (svc *MailService) send(ctx context.Context, email, title, body string, contentType mail.ContentType) error {
	var err error
	msg := mail.NewMsg()
	if err = msg.From(fmt.Sprintf("%s<%s>", svc.from, svc.username)); err != nil {
		return err
	}
	if err = msg.To(email); err != nil {
		return err
	}
	msg.Subject(title)
	msg.SetBodyString(contentType, body)
	if err = svc.client.DialAndSendWithContext(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (svc *MailService) SendString(ctx context.Context, email, title, body string) error {
	return svc.send(ctx, email, title, body, mail.TypeTextPlain)
}

func (svc *MailService) SendHTML(ctx context.Context, email, title, body string) error {
	return svc.send(ctx, email, title, body, mail.TypeTextHTML)
}

func (svc *MailService) Ping(ctx context.Context) error {
	return svc.send(ctx, svc.username, "hello", "server is starting", mail.TypeTextPlain)
}
