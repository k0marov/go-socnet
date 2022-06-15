package service

import (
	"core/core_values"
	"posts/domain/entities"
	"posts/domain/values"
)

type PostDeleter = func(values.PostId) error
type PostLikeToggler = func(values.PostId, core_values.UserId) error
type PostCreater = func(values.NewPostData) error
type PostsGetter = func(core_values.UserId) ([]entities.Post, error)
