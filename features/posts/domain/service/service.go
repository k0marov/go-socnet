package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/likeable"
	"github.com/k0marov/socnet/features/posts/domain/validators"
	profile_service "github.com/k0marov/socnet/features/profiles/domain/service"
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

func NewPostsGetter(getProfile profile_service.ProfileGetter, getPosts store.PostsGetter, checkLiked likeable.LikeChecker) PostsGetter {
	return func(authorId, caller core_values.UserId) (ctxPosts []entities.ContextedPost, err error) {
		author, err := getProfile(authorId, caller)
		if err == core_errors.ErrNotFound {
			return []entities.ContextedPost{}, client_errors.NotFound
		}
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while getting the ctxPosts author: %w", err)
		}
		posts, err := getPosts(authorId)
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while getting ctxPosts from store: %w", err)
		}
		for _, post := range posts {
			isLiked, err := checkLiked(post.Id, caller)
			if err != nil {
				return []entities.ContextedPost{}, fmt.Errorf("while checking if post is liked: %w", err)
			}
			ctxPost := entities.ContextedPost{
				Id:        post.Id,
				Author:    author,
				Text:      post.Text,
				Images:    post.Images,
				CreatedAt: post.CreatedAt,
				Likes:     post.Likes,
				IsLiked:   isLiked,
				IsMine:    authorId == caller,
			}
			ctxPosts = append(ctxPosts, ctxPost)
		}
		return
	}
}
