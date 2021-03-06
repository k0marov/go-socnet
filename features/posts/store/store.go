package store

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/deletable"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/store"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	"github.com/k0marov/go-socnet/features/posts/store/file_storage"
)

type (
	DBPostsGetter func(core_values.UserId) ([]models.PostModel, error)

	DBPostCreator     func(newPost models.PostToCreate) (values.PostId, error)
	DBPostImagesAdder func(values.PostId, []models.PostImageModel) error
)

// TODO: get rid of complexity by removing the "deleting on failure" logic by using transactions ?
func NewStorePostCreator(
	createPost DBPostCreator, storeImages file_storage.PostImageFilesCreator, addImages DBPostImagesAdder,
	deletePost deletable.ForceDeleter, deleteImages file_storage.PostFilesDeleter) store.PostCreator {
	return func(post values.NewPostData, createdAt time.Time) error {
		postToCreate := models.PostToCreate{
			Author:    post.Author,
			Text:      post.Text,
			CreatedAt: createdAt,
		}
		postId, err := createPost(postToCreate)
		if err != nil {
			return core_err.Rethrow("creating a post in db", err)
		}
		imagePaths, err := storeImages(postId, post.Author, post.Images)
		if err != nil || len(imagePaths) != len(post.Images) {
			deletePost(postId)
			return core_err.Rethrow("storing image files", err)
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
			return core_err.Rethrow("adding image paths to db", err)
		}
		return nil
	}
}

func NewStorePostDeleter(deletePost deletable.ForceDeleter, deleteFiles file_storage.PostFilesDeleter) store.PostDeleter {
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
			return []entities.Post{}, core_err.Rethrow("getting posts from db", err)
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
