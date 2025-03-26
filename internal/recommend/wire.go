package recommend

import "github.com/KNICEX/InkFlow/internal/recommend/internal/service"

func InitSyncService() SyncService {
	return service.NewSyncService()
}
