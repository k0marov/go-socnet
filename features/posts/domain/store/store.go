package store

import (
	"core/core_values"
	"posts/domain/entities"
	"posts/domain/values"
)

type StorePostsGetter func(authorId core_values.UserId) ([]entities.Post, error)
type StoreAuthorGetter func(postId values.PostId) (core_values.UserId, error)
type StorePostDeleter func(postId values.PostId) error
