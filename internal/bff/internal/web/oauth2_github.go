package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	ijwt "github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
)

type GithubOAuth2Handler struct {
	svc     user.OAuth2Service[user.GithubInfo]
	userSvc user.Service
	OAuth2Handler[user.GithubInfo]
}

func NewGithubOAuth2Handler(svc user.OAuth2Service[user.GithubInfo], userSvc user.Service, jwtHandler ijwt.Handler, l logx.Logger) *GithubOAuth2Handler {
	return &GithubOAuth2Handler{
		svc: svc,

		OAuth2Handler: OAuth2Handler[user.GithubInfo]{
			svc:          svc,
			stateKey:     "github-oauth2-123",
			callBackPath: "/oauth2/github",
			getDomain: func(ctx context.Context, t user.GithubInfo) (user.User, error) {
				return userSvc.FindOrCreateByGithub(ctx, t)
			},
			l:       l,
			Handler: jwtHandler,
		},
	}
}

func (o *GithubOAuth2Handler) RegisterRoutes(g *gin.Engine) {
	g.GET("/authUrl", ginx.Wrap(o.l, o.AuthUrl))
	g.GET("/callback", ginx.WrapBody(o.l, o.Callback))
}
