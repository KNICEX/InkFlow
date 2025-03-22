package ink

import (
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/service"
)

type Author = domain.Author
type Status = domain.Status
type Tags = domain.Tags
type Category = domain.Category

type Ink = domain.Ink

type Service = service.InkService

var ErrNoPermission = service.ErrNoPermission

const (
	StatusPublished = domain.InkStatusPublished
	StatusRejected  = domain.InkStatusRejected
)
