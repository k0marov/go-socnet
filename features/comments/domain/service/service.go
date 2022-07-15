package service

import (
	"github.com/k0marov/go-socnet/core/abstract/deletable"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable"
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/contexters"
	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/store"
	"github.com/k0marov/go-socnet/features/comments/domain/validators"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
)

type (
	PostCommentsGetter func(post post_values.PostId, caller core_values.UserId) ([]entities.ContextedComment, error)
	CommentCreator     func(newComment values.NewCommentValue) (entities.ContextedComment, error)
	CommentLikeToggler func(values.CommentId, core_values.UserId) error
	CommentDeleter     func(comment values.CommentId, caller core_values.UserId) error
)

func NewPostCommentsGetter(getComments store.CommentsGetter, addContexts contexters.CommentListContextAdder) PostCommentsGetter {
	return func(post post_values.PostId, caller core_values.UserId) ([]entities.ContextedComment, error) {
		comments, err := getComments(post)
		if err != nil {
			return []entities.ContextedComment{}, core_err.Rethrow("getting post contextedComments from store", err)
		}
		contextedComments, err := addContexts(comments, caller)
		if err != nil {
			return []entities.ContextedComment{}, core_err.Rethrow("adding contexts to comments", err)
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
			return entities.ContextedComment{}, core_err.Rethrow("getting author's profile", err)
		}

		createdAt := time.Now().UTC()
		newId, err := createComment(newComment, createdAt)
		if err != nil {
			return entities.ContextedComment{}, core_err.Rethrow("creating new comment", err)
		}

		comment := entities.ContextedComment{
			Comment: entities.Comment{
				CommentModel: models.CommentModel{
					Id:        newId,
					AuthorId:  newComment.Author,
					Text:      newComment.Text,
					CreatedAt: createdAt.Unix(),
				},
				Likes: 0,
			},
			OwnLikeContext: likeable_contexters.OwnLikeContext{
				IsLiked: false,
				IsMine:  true,
			},
			Author: author,
		}
		return comment, nil
	}
}

func NewCommentLikeToggler(safeToggleLike ownable_likeable.SafeLikeToggler) CommentLikeToggler {
	return CommentLikeToggler(safeToggleLike)
}

func NewCommentDeleter(delete deletable.Deleter) CommentDeleter {
	return CommentDeleter(delete)
}
