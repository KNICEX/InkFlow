package interactive

import (
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/events"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/service"
)

type Interactive = domain.Interactive
type ReadRecord = domain.ViewRecord
type LikeRecord = domain.LikeRecord

type Service = service.InteractiveService
type InkViewConsumer = events.InkViewConsumer
