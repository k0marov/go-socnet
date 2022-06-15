package service

import (
	"core/core_values"
	"posts/domain/entities"
	"posts/domain/values"
)

type PostDeleter func(post values.PostId, fromUser core_values.UserId) error
type PostLikeToggler func(values.PostId, core_values.UserId) error
type PostCreator func(values.NewPostData) (entities.Post, error)
type PostsGetter func(core_values.UserId) ([]entities.Post, error)
