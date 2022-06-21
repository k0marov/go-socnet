package service

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	PostCommentsGetter func(post post_values.PostId) ([]entities.Comment, error)
	CommentCreator     func(newComment values.NewCommentValue) (entities.Comment, error)
	CommentLikeToggler func(values.CommentId, core_values.UserId) error
)

func NewPostCommentsGetter() PostCommentsGetter {
	return func(post post_values.PostId) ([]entities.Comment, error) {
		panic("unimplemented")
	}
}

func NewCommentCreator() CommentCreator {
	return func(newComment values.NewCommentValue) (entities.Comment, error) {
		panic("unimplemented")
	}
}

func NewCommentLikeToggler() CommentLikeToggler {
	return func(comment values.CommentId, caller core_values.UserId) error {
		panic("unimplemented")
	}
}
