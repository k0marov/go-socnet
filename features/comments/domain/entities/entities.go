package entities

import (
	"github.com/k0marov/socnet/features/comments/domain/values"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type ContextedComment struct {
	Id        values.CommentId
	Author    profile_entities.ContextedProfile
	Text      string
	CreatedAt time.Time
	Likes     int

	IsLiked bool
	IsMine  bool
}
