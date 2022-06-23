package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type Post struct {
	Id        values.PostId
	AuthorId  core_values.UserId
	Text      string
	Images    []values.PostImage
	CreatedAt time.Time
	Likes     int
}

type ContextedPost struct {
	Post
	Author  profile_entities.ContextedProfile
	IsLiked bool
	IsMine  bool
}
