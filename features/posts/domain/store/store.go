package store

import (
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"time"

	"github.com/k0marov/socnet/core/core_values"
)

type PostsGetter func(authorId core_values.UserId) ([]entities.Post, error)
type AuthorGetter func(postId values.PostId) (core_values.UserId, error)
type PostDeleter func(postId values.PostId, authorId core_values.UserId) error
type PostCreator func(post values.NewPostData, createdAt time.Time) error
