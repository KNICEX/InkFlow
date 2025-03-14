package user

import (
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/internal/user/internal/service/oauth2"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitUserService(db *gorm.DB, cmd redis.Cmdable, l logx.Logger) Service {
	//wire.Build(
	//	dao.NewGormUserDAO,
	//	repo.NewCachedUserRepo,
	//	service.NewUserService,
	//)
	node := snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	d := dao.NewGormUserDAO(db, node)
	c := cache.NewRedisUserCache(cmd)
	r := repo.NewCachedUserRepo(d, c, l)
	return service.NewUserService(r, l)
}

func InitGithubOAuth2Service(l logx.Logger) OAuth2Service[GithubInfo] {
	type Config struct {
		ClientId       string `mapstructure:"client_id"`
		ClientSecret   string `mapstructure:"client_secret"`
		RedirectDomain string `mapstructure:"redirect_domain"`
	}

	var cfg Config
	if err := viper.UnmarshalKey("oauth2.github", &cfg); err != nil {
		panic(err)
	}

	return oauth2.NewGithubService(cfg.ClientId, cfg.ClientSecret, cfg.RedirectDomain, l)
}
