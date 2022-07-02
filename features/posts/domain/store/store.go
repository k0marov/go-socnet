package store

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type PostsGetter func(authorId core_values.UserId) ([]entities.Post, error)
type PostDeleter func(postId values.PostId, authorId core_values.UserId) error
type PostCreator func(post values.NewPostData, createdAt time.Time) error
