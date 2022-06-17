package entities

import (
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type ContextedPost struct {
	Id        values.PostId
	Author    profile.ContextedProfile
	Text      string
	Images    []values.PostImage
	CreatedAt time.Time
	IsLiked   bool
	IsMine    bool
}
