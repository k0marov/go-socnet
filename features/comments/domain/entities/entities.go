package entities

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/features/comments/domain/models"
	profile_entities "github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

type Comment struct {
	models.CommentModel
	Likes int
}

type ContextedComment struct {
	Comment
	contexters.OwnLikeContext
	Author profile_entities.ContextedProfile
}
