package action

import "github.com/KNICEX/InkFlow/internal/action/internal/service"

func InitService() Service {
	return &service.DoNothingActionService{}
}
