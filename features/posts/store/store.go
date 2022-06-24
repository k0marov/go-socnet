package store

import (
	"fmt"
	"time"

	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/likeable"
	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/store"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	"github.com/k0marov/go-socnet/features/posts/store/file_storage"
)

type (
	DBPostsGetter  func(core_values.UserId) ([]models.PostModel, error)
	DBAuthorGetter func(values.PostId) (core_values.UserId, error)

	DBPostCreator     func(newPost models.PostToCreate) (values.PostId, error)
	DBPostImagesAdder func(values.PostId, []models.PostImageModel) error
	DBPostDeleter     func(values.PostId) error
)

func NewStorePostCreator(
	createPost DBPostCreator, storeImages file_storage.PostImageFilesCreator, addImages DBPostImagesAdder,
	deletePost DBPostDeleter, deleteImages file_storage.PostFilesDeleter) store.PostCreator {
	return func(post values.NewPostData, createdAt time.Time) error {
		postToCreate := models.PostToCreate{
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
		var postImages []models.PostImageModel
		for i, path := range imagePaths {
			postImages = append(postImages, models.PostImageModel{
				Path:  path,
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

func NewStorePostsGetter(getter DBPostsGetter, likesGetter likeable.LikesCountGetter) store.PostsGetter {
	return func(authorId core_values.UserId) (posts []entities.Post, err error) {
		models, err := getter(authorId)
		if err != nil {
			return []entities.Post{}, fmt.Errorf("while getting posts from db: %w", err)
		}
		for _, model := range models {
			likes, err := likesGetter(model.Id)
			if err != nil {
				return []entities.Post{}, fmt.Errorf("error while getting likes count of a post: %w", err)
			}
			post := entities.Post{
				PostModel: model,
				Images:    entities.ImagePathsToUrls(model.Images),
				Likes:     likes,
			}
			posts = append(posts, post)
		}
		return
	}
}

func NewStoreAuthorGetter(authorGetter DBAuthorGetter) store.AuthorGetter {
	return store.AuthorGetter(authorGetter)
}
