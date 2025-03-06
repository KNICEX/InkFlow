package user

import (
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/oauth2"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitUserService(db *gorm.DB, cmd redis.Cmdable, l logx.Logger) Service {
	wire.Build(
		dao.NewGormUserDAO,
		repo.NewCachedUserRepo,
		service.NewUserService,
	)
	return service.NewUserService(nil, nil)
}

func InitGithubOAuth2Service(l logx.Logger) OAuth2Service[GithubInfo] {
	type Config struct {
		ClientId       string `yaml:"client_id"`
		ClientSecret   string `yaml:"client_secret"`
		RedirectDomain string `yaml:"redirect_domain"`
	}

	var cfg Config
	if err := viper.UnmarshalKey("oauth2.github", &cfg); err != nil {
		panic(err)
	}

	return oauth2.NewGithubService(cfg.ClientId, cfg.ClientSecret, cfg.RedirectDomain, l)
}
