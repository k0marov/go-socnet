package models

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"time"
)

type PostToCreate struct {
	Author    core_values.UserId
	Text      string
	CreatedAt time.Time
}

type PostModel struct {
	Id        values.PostId
	Author    core_values.UserId
	Text      string
	CreatedAt time.Time
	Images    []values.PostImage
}
