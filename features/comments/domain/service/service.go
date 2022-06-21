package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/store"
	"github.com/k0marov/socnet/features/comments/domain/validators"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
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
			comments = append(comments, modelToEntity(model))
		}
		return comments, nil
	}
}

func NewCommentCreator(validate validators.CommentValidator, createComment store.Creator) CommentCreator {
	return func(newComment values.NewCommentValue) (entities.Comment, error) {
		clientErr, isValid := validate(newComment)
		if !isValid {
			return entities.Comment{}, clientErr
		}
		commentModel, err := createComment(newComment)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.Comment{}, client_errors.NotFound
			}
			return entities.Comment{}, fmt.Errorf("while creating new comment: %w", err)
		}
		comment := modelToEntity(commentModel)
		return comment, nil
	}
}

func NewCommentLikeToggler(checkLiked store.LikeChecker, like store.Liker, unlike store.Unliker) CommentLikeToggler {
	return func(comment values.CommentId, caller core_values.UserId) error {
		isLiked, err := checkLiked(comment, caller)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return client_errors.NotFound
			}
			return fmt.Errorf("while checking if comment is liked: %w", err)
		}
		if isLiked {
			err = unlike(comment, caller)
			if err != nil {
				return fmt.Errorf("while unliking a comment: %w", err)
			}
		} else {
			err = like(comment, caller)
			if err != nil {
				return fmt.Errorf("while liking a comment: %w", err)
			}
		}
		return nil
	}
}

func modelToEntity(model models.CommentModel) entities.Comment {
	return entities.Comment{Id: model.Id}
}
