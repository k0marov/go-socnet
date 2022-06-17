package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile "github.com/k0marov/socnet/features/profiles/domain/entities"
	"time"
)

type Post struct {
	Id        values.PostId
	Author    profile.ContextedProfile
	Text      string
	Images    []core_values.FileURL
	CreatedAt time.Time
}

//type ContextedPost struct {
//	Post
//	IsLiked bool
//	IsMine  bool
//}
