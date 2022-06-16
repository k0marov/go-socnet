package service_test

import (
	"testing"

	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/service"
	"github.com/k0marov/socnet/features/posts/domain/values"

	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
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
	post := RandomString()
	caller := RandomString()
	t.Run("checking author", func(t *testing.T) {

		t.Run("error case - liking your own post", func(t *testing.T) {
			getAuthor := func(postId values.PostId) (core_values.UserId, error) {
				if postId == post {
					return caller, nil
				}
				panic("unexpected args")
			}
			err := service.NewPostLikeToggler(getAuthor, nil, nil, nil)(post, caller)
			AssertError(t, err, client_errors.LikingYourself)
		})
		t.Run("error case - post does not exist", func(t *testing.T) {
			getAuthor := func(postId values.PostId) (core_values.UserId, error) {
				return "", core_errors.ErrNotFound
			}
			err := service.NewPostLikeToggler(getAuthor, nil, nil, nil)(post, caller)
			AssertError(t, err, client_errors.NotFound)
		})
		t.Run("error case - checking author throws", func(t *testing.T) {
			getAuthor := func(values.PostId) (core_values.UserId, error) {
				return "", RandomError()
			}
			err := service.NewPostLikeToggler(getAuthor, nil, nil, nil)(post, caller)
			AssertSomeError(t, err)
		})
	})

	t.Run("liking/unliking, since author is not caller", func(t *testing.T) {
		getAuthor := func(values.PostId) (core_values.UserId, error) {
			return RandomString(), nil
		}
		t.Run("error case - checking if post is liked throws", func(t *testing.T) {
			likeChecker := func(postId values.PostId, callerId core_values.UserId) (bool, error) {
				if postId == post && callerId == caller {
					return false, RandomError()
				}
				panic("unexpected args")
			}
			err := service.NewPostLikeToggler(getAuthor, likeChecker, nil, nil)(post, caller)
			AssertSomeError(t, err)
		})
		t.Run("post not liked - like it", func(t *testing.T) {
			likeChecker := func(postId values.PostId, callerId core_values.UserId) (bool, error) {
				return false, nil
			}
			t.Run("happy case", func(t *testing.T) {
				liker := func(postId values.PostId, liker core_values.UserId) error {
					if postId == post && liker == caller {
						return nil
					}
					panic("unexpected args")
				}
				err := service.NewPostLikeToggler(getAuthor, likeChecker, liker, nil)(post, caller)
				AssertNoError(t, err)
			})
			t.Run("error case - store throws", func(t *testing.T) {
				liker := func(values.PostId, core_values.UserId) error {
					return RandomError()
				}
				err := service.NewPostLikeToggler(getAuthor, likeChecker, liker, nil)(post, caller)
				AssertSomeError(t, err)
			})
		})
		t.Run("post already liked - unlike it", func(t *testing.T) {
			likeChecker := func(values.PostId, core_values.UserId) (bool, error) {
				return true, nil
			}
			t.Run("happy case", func(t *testing.T) {
				unliker := func(postId values.PostId, unliker core_values.UserId) error {
					if postId == post && unliker == caller {
						return nil
					}
					panic("unexpected args")
				}
				err := service.NewPostLikeToggler(getAuthor, likeChecker, nil, unliker)(post, caller)
				AssertNoError(t, err)
			})
			t.Run("error case - store throws", func(t *testing.T) {
				unliker := func(values.PostId, core_values.UserId) error {
					return RandomError()
				}
				err := service.NewPostLikeToggler(getAuthor, likeChecker, nil, unliker)(post, caller)
				AssertSomeError(t, err)
			})
		})
	})
}

func TestPostCreator(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {

	})
	t.Run("error case - validation fails", func(t *testing.T) {

	})
	t.Run("error case - store returns error", func(t *testing.T) {

	})
}
