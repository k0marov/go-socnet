package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/store"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	PostCommentsGetter func(post post_values.PostId) ([]entities.Comment, error)
	CommentCreator     func(newComment values.NewCommentValue) (entities.Comment, error)
	CommentLikeToggler func(values.CommentId, core_values.UserId) error
)

func NewPostCommentsGetter(getComments store.CommentsGetter) PostCommentsGetter {
	return func(post post_values.PostId) ([]entities.Comment, error) {
		models, err := getComments(post)
		if err != nil {
			return []entities.Comment{}, fmt.Errorf("while getting post comments from store: %w", err)
		}
		var comments []entities.Comment
		for _, model := range models {
			comment := entities.Comment{Id: model.Id}
			comments = append(comments, comment)
		}
		return comments, nil
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
