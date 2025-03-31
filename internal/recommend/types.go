package recommend

import (
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/event"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/service"
)

type SyncService service.SyncService
type Service = service.RecommendService

type SyncConsumer = event.SyncConsumer

type Ink = domain.Ink
