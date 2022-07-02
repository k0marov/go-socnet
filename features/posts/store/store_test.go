package store_test

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"reflect"
	"testing"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	"github.com/k0marov/go-socnet/features/posts/store"
)

func TestStorePostCreator(t *testing.T) {
	tNewPost := RandomNewPostData()
	postId := RandomString()
	createdAt := time.Now()
	var imagePaths []core_values.StaticPath
	var wantPostImages []models.PostImageModel
	for _, img := range tNewPost.Images {
		path := RandomString()
		imagePaths = append(imagePaths, path)
		wantPostImages = append(wantPostImages, models.PostImageModel{Path: path, Index: img.Index})
	}

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
	addImages := func(post values.PostId, images []models.PostImageModel) error {
		if post == postId && reflect.DeepEqual(images, wantPostImages) {
			return nil
		}
		panic("unimplemented")
	}
	t.Run("error case - addImages returns an error", func(t *testing.T) {
		addImages := func(values.PostId, []models.PostImageModel) error {
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

func TestStorePostsGetter(t *testing.T) {
	author := RandomId()
	postModels := []models.PostModel{RandomPostModel()}
	likes := RandomInt()
	dbGetter := func(authorId core_values.UserId) ([]models.PostModel, error) {
		if authorId == author {
			return postModels, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting posts from db throws", func(t *testing.T) {
		dbGetter := func(core_values.UserId) ([]models.PostModel, error) {
			return nil, RandomError()
		}
		_, err := store.NewStorePostsGetter(dbGetter, nil)(author)
		AssertSomeError(t, err)
	})
	likesGetter := func(targetId string) (int, error) {
		if targetId == postModels[0].Id {
			return likes, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting likes throws", func(t *testing.T) {
		likesGetter := func(string) (int, error) {
			return 0, RandomError()
		}
		_, err := store.NewStorePostsGetter(dbGetter, likesGetter)(author)
		AssertSomeError(t, err)
	})
	gotPosts, err := store.NewStorePostsGetter(dbGetter, likesGetter)(author)
	AssertNoError(t, err)
	wantPosts := []entities.Post{{
		PostModel: postModels[0],
		Images:    entities.ImagePathsToUrls(postModels[0].Images),
		Likes:     likes,
	}}
	Assert(t, gotPosts, wantPosts, "returned posts")
}
