package store_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store"
	"github.com/k0marov/socnet/features/posts/store/models"
	"reflect"
	"testing"
	"time"
)

func TestStorePostCreator(t *testing.T) {
	tNewPost := RandomNewPostData()
	postId := RandomString()
	createdAt := time.Now()
	imagePaths := []core_values.StaticFilePath{RandomString(), RandomString()}

	createPost := func(newPost models.PostToCreate) (values.PostId, error) {
		if newPost.Author == tNewPost.Author && newPost.Text == tNewPost.Text && TimeAlmostEqual(newPost.CreatedAt, createdAt) {
			return postId, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - createPost returns an error", func(t *testing.T) {
		createPost := func(models.PostToCreate) (values.PostId, error) {
			return "", RandomError()
		}
		sut := store.NewStorePostCreator(createPost, nil, nil)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
	})
	storeImages := func(post values.PostId, author core_values.UserId, images []core_values.FileData) ([]core_values.StaticFilePath, error) {
		if post == postId && author == tNewPost.Author && reflect.DeepEqual(images, tNewPost.Images) {
			return imagePaths, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - storeImages returns an error", func(t *testing.T) {
		storeImages := func(values.PostId, core_values.UserId, []core_values.FileData) ([]core_values.StaticFilePath, error) {
			return []core_values.StaticFilePath{}, RandomError()
		}
		sut := store.NewStorePostCreator(createPost, storeImages, nil)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
	})
	addImages := func(post values.PostId, images []core_values.StaticFilePath) error {
		if post == postId && reflect.DeepEqual(images, imagePaths) {
			return nil
		}
		panic("unimplemented")
	}
	t.Run("error case - addImages returns an error", func(t *testing.T) {
		addImages := func(values.PostId, []core_values.StaticFilePath) error {
			return RandomError()
		}
		sut := store.NewStorePostCreator(createPost, storeImages, addImages)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
	})
	sut := store.NewStorePostCreator(createPost, storeImages, addImages)
	err := sut(tNewPost, createdAt)
	AssertNoError(t, err)
}

func TestStorePostDeleter(t *testing.T) {
	post := RandomString()
	author := RandomString()
	deletePost := func(postId values.PostId) error {
		if postId == post {
			return nil
		}
		panic("unexpected args")
	}
	t.Run("error case - delete post returns an error", func(t *testing.T) {
		deletePost := func(values.PostId) error {
			return RandomError()
		}
		sut := store.NewStorePostDeleter(deletePost, nil)
		err := sut(post, author)
		AssertSomeError(t, err)
	})
	deleteFiles := func(postId values.PostId, userId core_values.UserId) error {
		if postId == post && userId == author {
			return nil
		}
		panic("unexpected args")
	}
	t.Run("error case - delete files returns an error", func(t *testing.T) {
		deleteFiles := func(values.PostId, core_values.UserId) error {
			return RandomError()
		}
		sut := store.NewStorePostDeleter(deletePost, deleteFiles)
		err := sut(post, author)
		AssertSomeError(t, err)
	})

	sut := store.NewStorePostDeleter(deletePost, deleteFiles)
	err := sut(post, author)
	AssertNoError(t, err)
}
