package review

import (
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
)

type Service = service.Service

type AsyncService = service.AsyncService

type FailoverService = service.FailoverService

type Consumer = event.ReviewConsumer

type Ink = domain.Ink

type Result = domain.ReviewResult
