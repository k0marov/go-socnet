package entities

import (
	"github.com/k0marov/socnet/core/likeable/contexters"
	"github.com/k0marov/socnet/features/posts/domain/models"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
)

type Post struct {
	models.PostModel
	Likes int
}

type ContextedPost struct {
	Post
	contexters.LikeableContext
	Author profile_entities.ContextedProfile
}
