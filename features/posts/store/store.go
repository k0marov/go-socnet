package store

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/store"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/file_storage"
	"github.com/k0marov/socnet/features/posts/store/post_models"
	"time"
)

type (
	DBPostsGetter  func(core_values.UserId) ([]post_models.PostModel, error)
	DBLiker        func(values.PostId, core_values.UserId) error
	DBUnliker      func(values.PostId, core_values.UserId) error
	DBLikeChecker  func(values.PostId, core_values.UserId) (bool, error)
	DBAuthorGetter func(values.PostId) (core_values.UserId, error)

	DBPostCreator     func(newPost post_models.PostToCreate) (values.PostId, error)
	DBPostImagesAdder func(values.PostId, []values.PostImage) error
	DBPostDeleter     func(values.PostId) error
)

func NewStorePostCreator(
	createPost DBPostCreator, storeImages file_storage.PostImageFilesCreator, addImages DBPostImagesAdder,
	deletePost DBPostDeleter, deleteImages file_storage.PostFilesDeleter) store.PostCreator {
	return func(post values.NewPostData, createdAt time.Time) error {
		postToCreate := post_models.PostToCreate{
			Author:    post.Author,
			Text:      post.Text,
			CreatedAt: createdAt,
		}
		postId, err := createPost(postToCreate)
		if err != nil {
			return fmt.Errorf("while creating a post in db: %w", err)
		}
		imagePaths, err := storeImages(postId, post.Author, post.Images)
		if err != nil || len(imagePaths) != len(post.Images) {
			deletePost(postId)
			return fmt.Errorf("while storing image files: %w", err)
		}
		var postImages []values.PostImage
		for i, path := range imagePaths {
			postImages = append(postImages, values.PostImage{
				URL:   path,
				Index: post.Images[i].Index,
			})
		}
		err = addImages(postId, postImages)
		if err != nil {
			deletePost(postId)
			deleteImages(postId, post.Author)
			return fmt.Errorf("while adding image paths to db: %w", err)
		}
		return nil
	}
}

func NewStorePostDeleter(deletePost DBPostDeleter, deleteFiles file_storage.PostFilesDeleter) store.PostDeleter {
	return func(post values.PostId, author core_values.UserId) error {
		err := deletePost(post)
		if err != nil {
			return fmt.Errorf("error while deleting post from db: %w", err)
		}
		err = deleteFiles(post, author)
		if err != nil {
			return fmt.Errorf("error while deleting post files: %w", err)
		}
		return nil
	}
}

func NewStorePostsGetter(getter DBPostsGetter) store.PostsGetter {
	return store.PostsGetter(getter)
}

func NewStoreLiker(liker DBLiker) store.Liker {
	return store.Liker(liker)
}

func NewStoreLikeChecker(likeChecker DBLikeChecker) store.LikeChecker {
	return store.LikeChecker(likeChecker)
}
func NewStoreUnliker(unliker DBUnliker) store.Unliker {
	return store.Unliker(unliker)
}

func NewStoreAuthorGetter(authorGetter DBAuthorGetter) store.AuthorGetter {
	return store.AuthorGetter(authorGetter)
}
