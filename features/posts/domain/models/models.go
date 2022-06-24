package models

import (
	"time"

	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type PostToCreate struct {
	Author    core_values.UserId
	Text      string
	CreatedAt time.Time
}

type PostImageModel struct {
	Index int
	Path  core_values.StaticPath
}

type PostModel struct {
	Id        values.PostId
	AuthorId  core_values.UserId
	Text      string
	CreatedAt time.Time
	Images    []PostImageModel
}
