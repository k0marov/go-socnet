package service

import (
	"core/core_values"
	"fmt"
	"posts/domain/entities"

	"posts/domain/store"
	"posts/domain/values"
)

type (
	PostDeleter     func(post values.PostId, fromUser core_values.UserId) error
	PostLikeToggler func(values.PostId, core_values.UserId) error
	PostCreator     func(values.NewPostData) (entities.Post, error)
	PostsGetter     func(authorId core_values.UserId) ([]entities.Post, error)
)

func NewPostDeleter(getAuthor store.StoreAuthorGetter, deletePost store.StorePostDeleter) PostDeleter {
	return func(post values.PostId, fromUser core_values.UserId) error {
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

func NewPostsGetter(getPosts store.StorePostsGetter) PostsGetter {
	return func(authorId core_values.UserId) ([]entities.Post, error) {
		posts, err := getPosts(authorId)
		if err != nil {
			return []entities.Post{}, fmt.Errorf("while getting posts from store: %w", err)
		}
		return posts, nil
	}
}
