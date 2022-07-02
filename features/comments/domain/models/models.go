package models

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/values"
)

type CommentModel struct {
	Id        values.CommentId
	AuthorId  core_values.UserId
	Text      string
	CreatedAt time.Time
}
