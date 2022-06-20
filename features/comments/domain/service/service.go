package service

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	PostCommentsGetter func(post post_values.PostId) ([]entities.Comment, error)
	CommentAdder       func(post post_values.PostId, newComment values.NewCommentValue) (entities.Comment, error)
	CommentLikeToggler func(values.CommentId, core_values.UserId) error
)
