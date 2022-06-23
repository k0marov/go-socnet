package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type Post struct {
	Id        values.PostId
	Author    core_values.UserId
	Text      string
	Images    []values.PostImage
	CreatedAt time.Time
	Likes     int
}

type ContextedPost struct {
	Id        values.PostId
	Author    profile.ContextedProfile
	Text      string
	Images    []values.PostImage
	CreatedAt time.Time
	Likes     int

	IsLiked bool
	IsMine  bool
}
