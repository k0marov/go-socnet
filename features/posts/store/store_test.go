package store_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store"
	"github.com/k0marov/socnet/features/posts/store/post_models"
	"reflect"
	"testing"
	"time"
)

func TestStorePostCreator(t *testing.T) {
	tNewPost := RandomNewPostData()
	postId := RandomString()
	createdAt := time.Now()
	var imagePaths []core_values.StaticPath
	var wantPostImages []values.PostImage
	for _, img := range tNewPost.Images {
		path := RandomString()
		imagePaths = append(imagePaths, path)
		wantPostImages = append(wantPostImages, values.PostImage{URL: path, Index: img.Index})
	}

	createPost := func(newPost post_models.PostToCreate) (values.PostId, error) {
		if newPost.Author == tNewPost.Author && newPost.Text == tNewPost.Text && TimeAlmostEqual(newPost.CreatedAt, createdAt) {
			return postId, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - createPost returns an error", func(t *testing.T) {
		createPost := func(post_models.PostToCreate) (values.PostId, error) {
			return "", RandomError()
		}
		sut := store.NewStorePostCreator(createPost, nil, nil, nil, nil)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
	})
	storeImages := func(post values.PostId, author core_values.UserId, images []values.PostImageFile) ([]core_values.StaticPath, error) {
		if post == postId && author == tNewPost.Author && reflect.DeepEqual(images, tNewPost.Images) {
			return imagePaths, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - storeImages returns an error", func(t *testing.T) {
		storeImages := func(values.PostId, core_values.UserId, []values.PostImageFile) ([]core_values.StaticPath, error) {
			return []core_values.StaticPath{}, RandomError()
		}
		postDeleted := false
		deletePost := func(post values.PostId) error {
			if post == postId {
				postDeleted = true
				return nil
			}
			panic("unexpected args")
		}
		sut := store.NewStorePostCreator(createPost, storeImages, nil, deletePost, nil)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
		Assert(t, postDeleted, true, "post was deleted")
	})
	addImages := func(post values.PostId, images []values.PostImage) error {
		if post == postId && reflect.DeepEqual(images, wantPostImages) {
			return nil
		}
		panic("unimplemented")
	}
	t.Run("error case - addImages returns an error", func(t *testing.T) {
		addImages := func(values.PostId, []values.PostImage) error {
			return RandomError()
		}
		postDeleted := false
		imagesDeleted := false
		deletePost := func(post values.PostId) error {
			if post == postId {
				postDeleted = true
				return nil
			}
			panic("unexpected args")
		}
		deleteImages := func(post values.PostId, author core_values.UserId) error {
			if post == postId && author == tNewPost.Author {
				imagesDeleted = true
				return nil
			}
			panic("unexpected args")
		}
		sut := store.NewStorePostCreator(createPost, storeImages, addImages, deletePost, deleteImages)
		err := sut(tNewPost, createdAt)
		AssertSomeError(t, err)
		Assert(t, postDeleted, true, "post was deleted")
		Assert(t, imagesDeleted, true, "images were deleted")
	})
	t.Run("happy case", func(t *testing.T) {
		sut := store.NewStorePostCreator(createPost, storeImages, addImages, nil, nil)
		err := sut(tNewPost, createdAt)
		AssertNoError(t, err)
	})
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
