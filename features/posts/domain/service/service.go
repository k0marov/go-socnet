package service

import (
	"core/core_values"
	"posts/domain/entities"
	"posts/domain/values"
)

type (
	PostDeleter     func(post values.PostId, fromUser core_values.UserId) error
	PostLikeToggler func(values.PostId, core_values.UserId) error
	PostCreator     func(values.NewPostData) (entities.Post, error)
	PostsGetter     func(authorId core_values.UserId) ([]entities.Post, error)
)

func NewPostDeleter() PostDeleter {
	return func(post, fromUser core_values.UserId) error {
		panic("unimplemented")
	}
}

func NewPostLikeToggler() PostLikeToggler {
	return func(postId values.PostId, userId core_values.UserId) error {
		panic("unimplemented")
	}
}

func NewPostCreator() PostCreator {
	return func(newPost values.NewPostData) (entities.Post, error) {
		panic("unimplemented")
	}
}

func NewPostsGetter() PostsGetter {
	return func(authorId core_values.UserId) ([]entities.Post, error) {
		panic("unimplemented")
	}
}
