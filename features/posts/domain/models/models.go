package models

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type PostToCreate struct {
	Author    core_values.UserId
	Text      string
	CreatedAt time.Time
}

type PostImageModel struct {
	Index int                    `db:"ind"`
	Path  core_values.StaticPath `db:"path"`
}

type PostModel struct {
	Id        values.PostId      `db:"id"`
	AuthorId  core_values.UserId `db:"owner_id"`
	Text      string             `db:"textContent"`
	CreatedAt int64              `db:"createdAt"`
	Images    []PostImageModel
}
