package ink

import (
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/service"
)

type Author = domain.Author
type Status = domain.Status
type Category = domain.Category

type Ink = domain.Ink
type TagStats = domain.TagStats

type Service = service.InkService
type RankingService = service.RankingService

var ErrNoPermission = service.ErrNoPermission

const (
	StatusPending   = domain.InkStatusPending
	StatusPublished = domain.InkStatusPublished
	StatusRejected  = domain.InkStatusRejected
)
