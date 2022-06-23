package store

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/models"
	"github.com/k0marov/socnet/features/comments/domain/store"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"time"
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
			return []entities.Comment{}, fmt.Errorf("while getting post comments from db: %w", err)
		}
		for _, model := range commentModels {
			likes, err := getLikes(model.Id)
			if err != nil {
				return []entities.Comment{}, fmt.Errorf("while getting likes count for comment: %w", err)
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

func NewAuthorGetter(getAuthor DBAuthorGetter) store.AuthorGetter {
	return store.AuthorGetter(getAuthor)
}

func NewCommentCreator(createComment DBCommentCreator) store.Creator {
	return store.Creator(createComment)
}
