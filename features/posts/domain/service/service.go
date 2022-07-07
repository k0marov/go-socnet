package service

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/contexters"
	"github.com/k0marov/go-socnet/features/posts/domain/validators"

	"github.com/k0marov/go-socnet/features/posts/domain/entities"

	"github.com/k0marov/go-socnet/features/posts/domain/store"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type (
	PostDeleter     func(post values.PostId, caller core_values.UserId) error
	PostLikeToggler func(values.PostId, core_values.UserId) error
	PostCreator     func(values.NewPostData) error
	PostsGetter     func(fromAuthor, caller core_values.UserId) ([]entities.ContextedPost, error)
)

func NewPostDeleter(getAuthor ownable.OwnerGetter, deletePost store.PostDeleter) PostDeleter {
	return func(post values.PostId, caller core_values.UserId) error {
		author, err := getAuthor(post)
		if err != nil {
			return core_err.Rethrow("getting post author", err)
		}
		if author != caller {
			return client_errors.InsufficientPermissions
		}
		err = deletePost(post, author)
		if err != nil {
			return core_err.Rethrow("deleting post", err)
		}
		return nil
	}
}

func NewPostLikeToggler(safeToggleLike ownable_likeable.SafeLikeToggler) PostLikeToggler {
	return PostLikeToggler(safeToggleLike)
}

func NewPostCreator(validate validators.PostValidator, createPost store.PostCreator) PostCreator {
	return func(newPost values.NewPostData) error {
		clientError, ok := validate(newPost)
		if !ok {
			return clientError
		}
		err := createPost(newPost, time.Now())
		if err != nil {
			return core_err.Rethrow("creating a post in store", err)
		}
		return nil
	}
}

func NewPostsGetter(getPosts store.PostsGetter, addContext contexters.PostListContextAdder) PostsGetter {
	return func(authorId, caller core_values.UserId) ([]entities.ContextedPost, error) {
		posts, err := getPosts(authorId)
		if err != nil {
			return []entities.ContextedPost{}, core_err.Rethrow("getting posts from store", err)
		}
		ctxPosts, err := addContext(posts, caller)
		if err != nil {
			return []entities.ContextedPost{}, core_err.Rethrow("adding context to posts", err)
		}
		return ctxPosts, nil
	}
}
