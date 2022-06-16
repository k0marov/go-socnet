package store

import (
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"

	"github.com/k0marov/socnet/core/core_values"
)

type StorePostsGetter func(authorId core_values.UserId) ([]entities.Post, error)
type StoreAuthorGetter func(postId values.PostId) (core_values.UserId, error)
type StorePostDeleter func(postId values.PostId) error
type StoreLikeChecker func(postId values.PostId, liker core_values.UserId) (bool, error)
type StoreLiker func(postId values.PostId, liker core_values.UserId) error
type StoreUnliker func(postId values.PostId, unliker core_values.UserId) error
