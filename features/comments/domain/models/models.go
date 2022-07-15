package models

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
)

type CommentModel struct {
	Id        values.CommentId   `db:"id"`
	AuthorId  core_values.UserId `db:"owner_id"`
	Text      string             `db:"textContent"`
	CreatedAt int64              `db:"createdAt"`
}
