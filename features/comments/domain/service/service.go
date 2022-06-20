package service

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	GetPostComments func(post post_values.PostId) ([]entities.Comment, error)
	AddComment      func(post post_values.PostId, newComment values.NewCommentValue) (entities.Comment, error)
	ToggleLike      func(values.CommentId, core_values.UserId) error
)
