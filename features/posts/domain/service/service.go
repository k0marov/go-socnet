package service

import (
	"fmt"
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

func NewPostLikeToggler(getAuthor store.AuthorGetter, isLiked store.LikeChecker, like store.Liker, unlike store.Unliker) PostLikeToggler {
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

func NewPostsGetter(getProfile profile_service.ProfileGetter, getPosts store.PostsGetter, checkLiked store.LikeChecker) PostsGetter {
	return func(authorId, caller core_values.UserId) (posts []entities.ContextedPost, err error) {
		author, err := getProfile(authorId, caller)
		if err == core_errors.ErrNotFound {
			return []entities.ContextedPost{}, client_errors.NotFound
		}
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while getting the posts author: %w", err)
		}
		postModels, err := getPosts(authorId)
		if err != nil {
			return []entities.ContextedPost{}, fmt.Errorf("while getting posts from store: %w", err)
		}
		for _, postModel := range postModels {
			isLiked, err := checkLiked(postModel.Id, caller)
			if err != nil {
				return []entities.ContextedPost{}, fmt.Errorf("while checking if one of the posts is liked: %w", err)
			}
			post := entities.ContextedPost{
				Id:        postModel.Id,
				Author:    author,
				Text:      postModel.Text,
				Images:    postModel.Images,
				CreatedAt: postModel.CreatedAt,
				Likes:     postModel.Likes,
				IsLiked:   isLiked,
				IsMine:    authorId == caller,
			}
			posts = append(posts, post)
		}
		return
	}
}
