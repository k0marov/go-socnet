package entities

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable/contexters"
	"github.com/k0marov/go-socnet/features/comments/domain/models"
	profile_entities "github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

type Comment struct {
	models.CommentModel
	Likes int
}

type ContextedComment struct {
	Comment
	contexters.LikeableContext
	Author profile_entities.ContextedProfile
}
