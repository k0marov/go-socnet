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
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	profile_service "github.com/k0marov/socnet/features/profiles/domain/service"
	"time"
)

type (
	PostCommentsGetter func(post post_values.PostId, caller core_values.UserId) ([]entities.ContextedComment, error)
	CommentCreator     func(newComment values.NewCommentValue) (entities.ContextedComment, error)
	CommentLikeToggler func(values.CommentId, core_values.UserId) error
)

func NewPostCommentsGetter(getComments store.CommentsGetter, getProfile profile_service.ProfileGetter, checkLiked store.LikeChecker) PostCommentsGetter {
	return func(post post_values.PostId, caller core_values.UserId) ([]entities.ContextedComment, error) {
		models, err := getComments(post)
		if err != nil {
			return []entities.ContextedComment{}, fmt.Errorf("while getting post comments from store: %w", err)
		}
		var comments []entities.ContextedComment
		for _, model := range models {
			author, err := getProfile(model.Author, caller)
			if err != nil {
				return []entities.ContextedComment{}, fmt.Errorf("while getting comment's author profile: %w", err)
			}
			isLiked, err := checkLiked(model.Id, caller)
			if err != nil {
				return []entities.ContextedComment{}, fmt.Errorf("while checking if comment is liked: %w", err)
			}
			comment := entities.ContextedComment{
				Id:        model.Id,
				Author:    author,
				Text:      model.Text,
				CreatedAt: model.CreatedAt,
				Likes:     model.Likes,

				IsLiked: isLiked,
				IsMine:  author.Id == caller,
			}
			comments = append(comments, comment)
		}
		return comments, nil
	}
}

func NewCommentCreator(validate validators.CommentValidator, getProfile profile_service.ProfileGetter, createComment store.Creator) CommentCreator {
	return func(newComment values.NewCommentValue) (entities.ContextedComment, error) {
		clientErr, isValid := validate(newComment)
		if !isValid {
			return entities.ContextedComment{}, clientErr
		}

		author, err := getProfile(newComment.Author, newComment.Author)
		if err != nil {
			return entities.ContextedComment{}, fmt.Errorf("while getting author's profile: %w", err)
		}

		newId, err := createComment(newComment, time.Now())
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.ContextedComment{}, client_errors.NotFound
			}
			return entities.ContextedComment{}, fmt.Errorf("while creating new comment: %w", err)
		}

		comment := entities.ContextedComment{
			Id:        newId,
			Author:    author,
			Text:      newComment.Text,
			CreatedAt: time.Now(),
			Likes:     0,
			IsLiked:   false,
			IsMine:    true,
		}
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
