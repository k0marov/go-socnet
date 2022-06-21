package store

import (
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/post_models"
	"time"

	"github.com/k0marov/socnet/core/core_values"
)

type PostsGetter func(authorId core_values.UserId) ([]post_models.PostModel, error)
type AuthorGetter func(postId values.PostId) (core_values.UserId, error)
type PostDeleter func(postId values.PostId, authorId core_values.UserId) error
type LikeChecker func(postId values.PostId, liker core_values.UserId) (bool, error)
type Liker func(postId values.PostId, liker core_values.UserId) error
type Unliker func(postId values.PostId, unliker core_values.UserId) error
type PostCreator func(post values.NewPostData, createdAt time.Time) error
