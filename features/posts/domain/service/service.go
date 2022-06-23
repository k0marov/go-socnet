package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/likeable"
	"github.com/k0marov/socnet/features/posts/domain/contexters"
	"github.com/k0marov/socnet/features/posts/domain/validators"
	"time"

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
	PostCreator     func(values.NewPostData) error
	PostsGetter     func(fromAuthor, caller core_values.UserId) ([]entities.ContextedPost, error)
)

func NewPostDeleter(getAuthor store.AuthorGetter, deletePost store.PostDeleter) PostDeleter {
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
		err = deletePost(post, author)
		if err != nil {
			return fmt.Errorf("while deleting post: %w", err)
		}
		return nil
	}
}

func NewPostLikeToggler(getAuthor store.AuthorGetter, toggleLike likeable.LikeToggler) PostLikeToggler {
	return func(postId values.PostId, caller core_values.UserId) error {
		author, err := getAuthor(postId)
		if err == core_errors.ErrNotFound {
			return client_errors.NotFound
		}
		if err != nil {
			return fmt.Errorf("while getting post author: %w", err)
		}

		err = toggleLike(postId, author, caller)
		if err != nil {
			return err
		}

		return nil
	}
}

func NewPostCreator(validate validators.PostValidator, createPost store.PostCreator) PostCreator {
	return func(newPost values.NewPostData) error {
		clientError, ok := validate(newPost)
		if !ok {
			return clientError
		}
		err := createPost(newPost, time.Now())
		if err != nil {
			return fmt.Errorf("while creating a post in store: %w", err)
		}
		return nil
	}
}

func NewPostsGetter(getPosts store.PostsGetter, addContext contexters.PostListContextAdder) PostsGetter {
	return func(authorId, caller core_values.UserId) ([]entities.ContextedPost, error) {
		posts, err := getPosts(authorId)
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while getting posts from store: %w", err)
		}
		ctxPosts, err := addContext(posts, caller)
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while adding context to posts: %w", err)
		}
		return ctxPosts, nil
	}
}
