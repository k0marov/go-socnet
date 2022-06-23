package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	"github.com/k0marov/socnet/features/comments/domain/contexters"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/models"
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

func NewPostCommentsGetter(getComments store.CommentsGetter, addContexts contexters.CommentListContextAdder) PostCommentsGetter {
	return func(post post_values.PostId, caller core_values.UserId) ([]entities.ContextedComment, error) {
		comments, err := getComments(post)
		if err != nil {
			return []entities.ContextedComment{}, fmt.Errorf("while getting post contextedComments from store: %w", err)
		}
		contextedComments, err := addContexts(comments, caller)
		if err != nil {
			return []entities.ContextedComment{}, fmt.Errorf("while adding contexts to comments: %w", err)
		}
		return contextedComments, nil
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

		createdAt := time.Now()
		newId, err := createComment(newComment, createdAt)
		if err != nil {
			return entities.ContextedComment{}, fmt.Errorf("while creating new comment: %w", err)
		}

		comment := entities.ContextedComment{
			Comment: entities.Comment{
				CommentModel: models.CommentModel{
					Id:        newId,
					AuthorId:  newComment.Author,
					Text:      newComment.Text,
					CreatedAt: createdAt,
				},
				Likes: 0,
			},
			LikeableContext: likeable_contexters.LikeableContext{
				IsLiked: false,
				IsMine:  true,
			},
			Author: author,
		}
		return comment, nil
	}
}

func NewCommentLikeToggler(getAuthor store.AuthorGetter, toggleLike likeable.LikeToggler) CommentLikeToggler {
	return func(comment values.CommentId, caller core_values.UserId) error {
		author, err := getAuthor(comment)
		if err == core_errors.ErrNotFound {
			return client_errors.NotFound
		}
		if err != nil {
			return fmt.Errorf("while checking comment author: %w", err)
		}
		err = toggleLike(comment, author, caller)
		if err != nil {
			return err
		}
		return nil
	}
}
