package comment

import (
	"github.com/KNICEX/InkFlow/internal/comment/internal/domain"
	"github.com/KNICEX/InkFlow/internal/comment/internal/service"
)

type Service = service.CommentService

type Comment = domain.Comment
type Payload = domain.Payload
type Commentator = domain.Commentator
type Stats = domain.CommentStats
