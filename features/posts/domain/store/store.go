package store

import (
	"core/core_values"
	"posts/domain/entities"
	"posts/domain/values"
)

type StorePostsGetter func(authorId core_values.UserId) ([]entities.Post, error)
type StoreAuthorGetter func(postId values.PostId) (core_values.UserId, error)
type StorePostDeleter func(postId values.PostId) error
type StoreLikeChecker func(postId values.PostId, liker core_values.UserId) (bool, error)
type StoreLiker func(postId values.PostId, liker core_values.UserId) error
type StoreUnliker func(postId values.PostId, unliker core_values.UserId) error
