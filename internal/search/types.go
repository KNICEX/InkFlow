package search

import (
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/event"
	"github.com/KNICEX/InkFlow/internal/search/internal/service"
)

type SyncService = service.SyncService
type Service = service.SearchService

type SyncConsumer = event.SyncConsumer

type Ink = domain.Ink
type Comment = domain.Comment
type User = domain.User
