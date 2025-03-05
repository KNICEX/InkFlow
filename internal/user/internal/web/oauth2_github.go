package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/oauth2"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	ijwt "github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
)

type GithubOAuth2Handler struct {
	svc     oauth2.Service[domain.GithubInfo]
	userSvc service.UserService
	oAuth2Handler[domain.GithubInfo]
}

func NewGithubOAuth2Handler(svc oauth2.Service[domain.GithubInfo], userSvc service.UserService, jwtHandler ijwt.Handler, l logx.Logger) *GithubOAuth2Handler {
	return &GithubOAuth2Handler{
		svc: svc,

		oAuth2Handler: oAuth2Handler[domain.GithubInfo]{
			svc:          svc,
			stateKey:     "github-oauth2-123",
			callBackPath: "/oauth2/github",
			getDomain: func(ctx context.Context, I domain.GithubInfo) (domain.User, error) {
				return userSvc.FindOrCreateByGithub(ctx, I)
			},
			logger:  l,
			Handler: jwtHandler,
		},
	}
}

func (o *GithubOAuth2Handler) RegisterRoutes(g *gin.Engine) {
	g.GET("/authUrl", ginx.Wrap(o.logger, o.AuthUrl))
	g.GET("/callback", ginx.WrapBody(o.logger, o.Callback))
}
