package code

import (
	"github.com/KNICEX/InkFlow/internal/code/internal/repo"
	"github.com/KNICEX/InkFlow/internal/code/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/code/internal/service"
	"github.com/KNICEX/InkFlow/internal/email"
	"github.com/redis/go-redis/v9"
)

type Service = service.Service

func InitEmailCodeService(cmd redis.Cmdable, emailSvc email.Service) Service {
	ca := cache.NewCodeCache(cmd)
	repository := repo.NewCodeRepo(ca)
	return service.NewCachedEmailCodeService(repository, emailSvc)
}
