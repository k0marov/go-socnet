package store

import (
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
)

type (
	CommentsGetter func(post post_values.PostId) ([]entities.Comment, error)
	Creator        func(newComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error)
)
