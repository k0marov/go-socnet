package store

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"time"
)

type (
	CommentsGetter func(post post_values.PostId) ([]entities.Comment, error)
	AuthorGetter   func(comment values.CommentId) (core_values.UserId, error)
	Creator        func(newComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error)
)
