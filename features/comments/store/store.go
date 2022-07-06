package store

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/store"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
)

type (
	DBCommentsGetter func(post post_values.PostId) ([]models.CommentModel, error)
	DBAuthorGetter   func(post post_values.PostId) (core_values.UserId, error)
	DBCommentCreator func(newComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error)
)

func NewCommentsGetter(getComments DBCommentsGetter, getLikes likeable.LikesCountGetter) store.CommentsGetter {
	return func(post post_values.PostId) (comments []entities.Comment, error error) {
		commentModels, err := getComments(post)
		if err != nil {
			return []entities.Comment{}, core_err.Rethrow("getting post comments from db", err)
		}
		for _, model := range commentModels {
			likes, err := getLikes(model.Id)
			if err != nil {
				return []entities.Comment{}, core_err.Rethrow("getting likes count for comment", err)
			}
			comment := entities.Comment{
				CommentModel: model,
				Likes:        likes,
			}
			comments = append(comments, comment)
		}
		return
	}
}

func NewCommentCreator(createComment DBCommentCreator) store.Creator {
	return store.Creator(createComment)
}
