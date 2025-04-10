package interactive

import (
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/events"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/service"
)

type Interactive = domain.Interactive
type Favorite = domain.Favorite
type ViewRecord = domain.ViewRecord
type LikeRecord = domain.LikeRecord
type UserStats = domain.UserStats
type FavoriteRecord = domain.FavoriteRecord

type Service = service.InteractiveService
type InkViewConsumer = events.InkViewConsumer
