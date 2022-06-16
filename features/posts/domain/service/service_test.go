package service_test

import (
	"core/client_errors"
	"core/core_errors"
	"core/core_values"
	. "core/test_helpers"
	"posts/domain/entities"
	"posts/domain/service"
	"posts/domain/values"
	"testing"
)

func TestPostsGetter(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		randomPosts := []entities.Post{RandomPost(), RandomPost(), RandomPost()}
		randomAuthor := RandomString()
		storePostsGetter := func(authorId core_values.UserId) ([]entities.Post, error) {
			if authorId == randomAuthor {
				return randomPosts, nil
			}
			panic("unexpected args")
		}
		sut := service.NewPostsGetter(storePostsGetter)
		gotPosts, err := sut(randomAuthor)
		AssertNoError(t, err)
		Assert(t, gotPosts, randomPosts, "the returned posts")
	})
	t.Run("error case - store throws an error", func(t *testing.T) {
		storeGetter := func(core_values.UserId) ([]entities.Post, error) {
			return []entities.Post{}, RandomError()
		}
		_, err := service.NewPostsGetter(storeGetter)("42")
		AssertSomeError(t, err)
	})
}

func TestPostDeleter(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		post := RandomString()
		caller := RandomString()
		postAuthor := caller
		authorGetter := func(postId values.PostId) (core_values.UserId, error) {
			if postId == post {
				return postAuthor, nil
			}
			panic("unexpected args")
		}
		deleted := false
		postDeleter := func(postId values.PostId) error {
			if postId == post {
				deleted = true
				return nil
			}
			panic("unexpected args")
		}
		sut := service.NewPostDeleter(authorGetter, postDeleter)
		err := sut(post, caller)
		AssertNoError(t, err)
		Assert(t, deleted, true, "post was deleted")
	})
	t.Run("error case - the calling user is not the post author", func(t *testing.T) {
		getAuthor := func(values.PostId) (core_values.UserId, error) {
			return RandomString(), nil
		}
		sut := service.NewPostDeleter(getAuthor, nil)
		err := sut(RandomString(), RandomString())
		AssertError(t, err, client_errors.InsufficientPermissions)
	})
	t.Run("error case - getting author returns post not found", func(t *testing.T) {
		getAuthor := func(values.PostId) (core_values.UserId, error) {
			return "", core_errors.ErrNotFound
		}
		sut := service.NewPostDeleter(getAuthor, nil)
		err := sut(RandomString(), RandomString())
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - getting author returns some other error", func(t *testing.T) {
		getAuthor := func(values.PostId) (core_values.UserId, error) {
			return "", RandomError()
		}
		err := service.NewPostDeleter(getAuthor, nil)(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("error case - deleting post returns error", func(t *testing.T) {
		caller := "42"
		getAuthor := func(values.PostId) (core_values.UserId, error) {
			return caller, nil
		}
		deletePost := func(values.PostId) error {
			return RandomError()
		}
		err := service.NewPostDeleter(getAuthor, deletePost)("33", caller)
		AssertSomeError(t, err)
	})
}

func TestPostLikeToggler(t *testing.T) {
	// post := RandomString()
	// caller := RandomString()
	t.Run("liking your own post", func(t *testing.T) {
		// sut := service.NewPostLikeToggler(nil, nil, nil)(post, caller)
	})
	t.Run("post not liked - like it", func(t *testing.T) {
		// likeChecker := func()
	})
	t.Run("post already liked - unlike it", func(t *testing.T) {

	})
}

func TestPostCreator(t *testing.T) {

}
