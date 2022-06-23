package models

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"time"
)

type CommentModel struct {
	Id        values.CommentId
	AuthorId  core_values.UserId
	Text      string
	CreatedAt time.Time
}
