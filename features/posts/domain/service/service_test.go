package service_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/entities"

	"github.com/k0marov/go-socnet/features/posts/domain/service"
	"github.com/k0marov/go-socnet/features/posts/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/core_errors"
	"github.com/k0marov/go-socnet/core/core_values"
	. "github.com/k0marov/go-socnet/core/test_helpers"
)

func TestPostsGetter(t *testing.T) {
	author := RandomId()
	caller := RandomString()
	posts := []entities.Post{RandomPost()}
	ctxPosts := []entities.ContextedPost{RandomContextedPost()}

	storePostsGetter := func(authorId core_values.UserId) ([]entities.Post, error) {
		if authorId == author {
			return posts, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - store throws an error", func(t *testing.T) {
		storeGetter := func(core_values.UserId) ([]entities.Post, error) {
			return []entities.Post{}, RandomError()
		}
		_, err := service.NewPostsGetter(storeGetter, nil)(author, caller)
		AssertSomeError(t, err)
	})
	contextAdder := func(postsList []entities.Post, callerId core_values.UserId) ([]entities.ContextedPost, error) {
		if reflect.DeepEqual(postsList, posts) && callerId == caller {
			return ctxPosts, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - context adder throws an error", func(t *testing.T) {
		contextAdder := func([]entities.Post, core_values.UserId) ([]entities.ContextedPost, error) {
			return nil, RandomError()
		}
		_, err := service.NewPostsGetter(storePostsGetter, contextAdder)(author, caller)
		AssertSomeError(t, err)
	})
	gotPosts, err := service.NewPostsGetter(storePostsGetter, contextAdder)(author, caller)
	AssertNoError(t, err)
	Assert(t, gotPosts, ctxPosts, "returned posts")
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
		postDeleter := func(postId values.PostId, authorId core_values.UserId) error {
			if postId == post && authorId == postAuthor {
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
		deletePost := func(values.PostId, core_values.UserId) error {
			return RandomError()
		}
		err := service.NewPostDeleter(getAuthor, deletePost)("33", caller)
		AssertSomeError(t, err)
	})
}

func TestPostLikeToggler(t *testing.T) {
	post := RandomString()
	author := RandomString()
	caller := RandomString()
	getAuthor := func(postId values.PostId) (core_values.UserId, error) {
		if postId == post {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting author throws", func(t *testing.T) {
		t.Run("post does not exist", func(t *testing.T) {
			getAuthor := func(postId values.PostId) (core_values.UserId, error) {
				return "", core_errors.ErrNotFound
			}
			err := service.NewPostLikeToggler(getAuthor, nil)(post, caller)
			AssertError(t, err, client_errors.NotFound)
		})
		t.Run("some other error", func(t *testing.T) {
			getAuthor := func(values.PostId) (core_values.UserId, error) {
				return "", RandomError()
			}
			err := service.NewPostLikeToggler(getAuthor, nil)(post, caller)
			AssertSomeError(t, err)
		})
	})
	toggleLike := func(target string, owner, callerId core_values.UserId) error {
		if target == post && owner == author && callerId == caller {
			return nil
		}
		panic("unexpected args")
	}
	t.Run("error case - toggling like throws, should FORWARD the error, since it can be a client error", func(t *testing.T) {
		wantErr := RandomError()
		toggleLike := func(target string, owner, callerId core_values.UserId) error {
			return wantErr
		}
		err := service.NewPostLikeToggler(getAuthor, toggleLike)(post, caller)
		AssertError(t, err, wantErr)
	})
	err := service.NewPostLikeToggler(getAuthor, toggleLike)(post, caller)
	AssertNoError(t, err)
}

func TestPostCreator(t *testing.T) {
	tNewPost := RandomNewPostData()
	t.Run("happy case", func(t *testing.T) {
		validator := func(newPost values.NewPostData) (client_errors.ClientError, bool) {
			if reflect.DeepEqual(newPost, tNewPost) {
				return client_errors.ClientError{}, true
			}
			panic("unexpected args")
		}
		storeCreator := func(newPost values.NewPostData, createdAt time.Time) error {
			if reflect.DeepEqual(newPost, tNewPost) && TimeAlmostNow(createdAt) {
				return nil
			}
			panic("unexpected args")
		}
		err := service.NewPostCreator(validator, storeCreator)(tNewPost)
		AssertNoError(t, err)
	})
	t.Run("error case - validation fails", func(t *testing.T) {
		wantErr := RandomClientError()
		validator := func(values.NewPostData) (client_errors.ClientError, bool) {
			return wantErr, false
		}
		err := service.NewPostCreator(validator, nil)(tNewPost)
		AssertError(t, err, wantErr)
	})
	t.Run("error case - store returns error", func(t *testing.T) {
		validator := func(values.NewPostData) (client_errors.ClientError, bool) {
			return client_errors.ClientError{}, true
		}
		storeCreator := func(values.NewPostData, time.Time) error {
			return RandomError()
		}
		err := service.NewPostCreator(validator, storeCreator)(tNewPost)
		AssertSomeError(t, err)
	})
}
