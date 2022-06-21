package service_test

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"testing"
)

func TestCommentCreator(t *testing.T) {

}

func TestCommentLikeToggler(t *testing.T) {
	comment := RandomString()
	caller := RandomString()
	t.Run("comment is already liked - unlike it", func(t *testing.T) {
		likeChecker := func(commentId values.CommentId, callerId core_values.UserId) (bool, error) {
			if comment == commentId && caller == callerId {
				return true, nil
			}
			panic("unexpected args")
		}
		unliker := func(commentId values.CommentId, unliker core_values.UserId) error {
			if commentId == comment && unliker == caller {
				return nil
			}
			panic("unexpected args")
		}
		t.Run("error case - unliker throws", func(t *testing.T) {
			unliker := func(values.CommentId, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewCommentLikeToggler(likeChecker, nil, unliker)(comment, caller)
			AssertSomeError(t, err)
		})
		err := service.NewCommentLikeToggler(likeChecker, nil, unliker)(comment, caller)
		AssertNoError(t, err)
	})
	t.Run("comment is not already liked - like it", func(t *testing.T) {
		likeChecker := func(commentId values.CommentId, callerId core_values.UserId) (bool, error) {
			if comment == commentId && caller == callerId {
				return false, nil
			}
			panic("unexpected args")
		}
		liker := func(commentId values.CommentId, liker core_values.UserId) error {
			if commentId == comment && liker == caller {
				return nil
			}
			panic("unexpected args")
		}
		t.Run("error case - liker throws", func(t *testing.T) {
			liker := func(values.CommentId, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewCommentLikeToggler(likeChecker, liker, nil)(comment, caller)
			AssertSomeError(t, err)
		})
		err := service.NewCommentLikeToggler(likeChecker, liker, nil)(comment, caller)
		AssertNoError(t, err)
	})
	t.Run("error case - like checker throws an error", func(t *testing.T) {
		t.Run("it is a not found error, should return client error", func(t *testing.T) {
			likeChecker := func(values.CommentId, core_values.UserId) (bool, error) {
				return false, core_errors.ErrNotFound
			}
			err := service.NewCommentLikeToggler(likeChecker, nil, nil)(comment, caller)
			AssertError(t, err, client_errors.NotFound)
		})
		t.Run("it is some other error", func(t *testing.T) {
			likeChecker := func(values.CommentId, core_values.UserId) (bool, error) {
				return false, RandomError()
			}
			err := service.NewCommentLikeToggler(likeChecker, nil, nil)(comment, caller)
			AssertSomeError(t, err)
		})
	})
}

func TestPostCommentsGetter(t *testing.T) {
	post := RandomString()
	commentModels := []models.CommentModel{RandomCommentModel(), RandomCommentModel()}
	wantComments := []entities.Comment{}
	for _, model := range commentModels {
		comment := entities.Comment{Id: model.Id}
		wantComments = append(wantComments, comment)
	}
	t.Run("happy case", func(t *testing.T) {
		getter := func(postId post_values.PostId) ([]models.CommentModel, error) {
			if postId == post {
				return commentModels, nil
			}
			panic("unexpected args")
		}
		got, err := service.NewPostCommentsGetter(getter)(post)
		AssertNoError(t, err)
		Assert(t, got, wantComments, "returned commentModels")
	})
	t.Run("error case - store returns some error", func(t *testing.T) {
		getter := func(post_values.PostId) ([]models.CommentModel, error) {
			return []models.CommentModel{}, RandomError()
		}
		_, err := service.NewPostCommentsGetter(getter)(post)
		AssertSomeError(t, err)
	})
}
