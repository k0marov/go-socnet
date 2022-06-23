package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/values"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type Comment struct {
	Id        values.CommentId
	AuthorId  core_values.UserId
	Text      string
	CreatedAt time.Time
	Likes     int
}

type ContextedComment struct {
	Comment
	Author          profile_entities.ContextedProfile
	IsLiked, IsMine bool
}
