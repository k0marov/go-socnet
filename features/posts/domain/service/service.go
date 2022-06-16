package service

import (
	"fmt"

	"github.com/k0marov/socnet/features/posts/domain/entities"

	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"

	"github.com/k0marov/socnet/features/posts/domain/store"
	"github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	PostDeleter     func(post values.PostId, caller core_values.UserId) error
	PostLikeToggler func(values.PostId, core_values.UserId) error
	PostCreator     func(values.NewPostData) (entities.Post, error)
	PostsGetter     func(authorId core_values.UserId) ([]entities.Post, error)
)

func NewPostDeleter(getAuthor store.StoreAuthorGetter, deletePost store.StorePostDeleter) PostDeleter {
	return func(post values.PostId, caller core_values.UserId) error {
		author, err := getAuthor(post)
		if err == core_errors.ErrNotFound {
			return client_errors.NotFound
		}
		if err != nil {
			return fmt.Errorf("while getting post author: %w", err)
		}
		if author != caller {
			return client_errors.InsufficientPermissions
		}
		err = deletePost(post)
		if err != nil {
			return fmt.Errorf("while deleting post: %w", err)
		}
		return nil
	}
}

func NewPostLikeToggler(getAuthor store.StoreAuthorGetter, isLiked store.StoreLikeChecker, like store.StoreLiker, unlike store.StoreUnliker) PostLikeToggler {
	return func(postId values.PostId, caller core_values.UserId) error {
		author, err := getAuthor(postId)
		if err == core_errors.ErrNotFound {
			return client_errors.NotFound
		}
		if err != nil {
			return fmt.Errorf("while getting post author: %w", err)
		}
		if author == caller {
			return client_errors.LikingYourself
		}

		alreadyLiked, err := isLiked(postId, caller)
		if err != nil {
			return fmt.Errorf("error while checking if post is liked: %w", err)
		}
		if alreadyLiked {
			err = unlike(postId, caller)
		} else {
			err = like(postId, caller)
		}
		if err != nil {
			return fmt.Errorf("while liking/unliking a post: %w", err)
		}
		return nil
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
