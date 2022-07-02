package service_test

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"reflect"
	"testing"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/service"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
	profile_entities "github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

func TestCommentCreator(t *testing.T) {
	newComment := RandomNewComment()
	createdId := RandomString()
	author := RandomContextedProfile()
	createdComment := entities.ContextedComment{
		Comment: entities.Comment{
			CommentModel: models.CommentModel{
				Id:        createdId,
				AuthorId:  newComment.Author,
				Text:      newComment.Text,
				CreatedAt: time.Now(),
			},
			Likes: 0,
		},
		LikeableContext: likeable_contexters.LikeableContext{
			IsLiked: false,
			IsMine:  true,
		},
		Author: author,
	}
	validator := func(gotComment values.NewCommentValue) (client_errors.ClientError, bool) {
		if gotComment == newComment {
			return client_errors.ClientError{}, true
		}
		panic("unexpected args")
	}
	t.Run("validator throws a client error", func(t *testing.T) {
		clientErr := RandomClientError()
		validator := func(value values.NewCommentValue) (client_errors.ClientError, bool) {
			return clientErr, false
		}
		_, err := service.NewCommentCreator(validator, nil, nil)(newComment)
		AssertError(t, err, clientErr)
	})
	profileGetter := func(target, caller core_values.UserId) (profile_entities.ContextedProfile, error) {
		if target == newComment.Author && caller == newComment.Author {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("profile getter throws", func(t *testing.T) {
		profileGetter := func(target, caller core_values.UserId) (profile_entities.ContextedProfile, error) {
			return profile_entities.ContextedProfile{}, RandomError()
		}
		_, err := service.NewCommentCreator(validator, profileGetter, nil)(newComment)
		AssertSomeError(t, err)
	})

	creator := func(gotComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error) {
		if gotComment == newComment && TimeAlmostNow(createdAt) {
			return createdId, nil
		}
		panic("unexpected args")
	}
	t.Run("creator throws", func(t *testing.T) {
		creator := func(values.NewCommentValue, time.Time) (values.CommentId, error) {
			return "", RandomError()
		}
		_, err := service.NewCommentCreator(validator, profileGetter, creator)(newComment)
		AssertSomeError(t, err)
	})
	sut := service.NewCommentCreator(validator, profileGetter, creator)
	gotCreated, err := sut(newComment)
	AssertNoError(t, err)
	Assert(t, TimeAlmostNow(gotCreated.CreatedAt), true, "createdAt is time.Now()")
	_, zoneOffset := gotCreated.CreatedAt.Zone()
	Assert(t, zoneOffset, 0, "time zone offset")
	gotCreated.CreatedAt = createdComment.CreatedAt
	Assert(t, gotCreated, createdComment, "the returned created comment")
}

func TestCommentLikeToggler(t *testing.T) {
	comment := RandomId()
	author := RandomId()
	caller := RandomId()
	authorGetter := func(commentId values.CommentId) (core_values.UserId, error) {
		if commentId == comment {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting author throws", func(t *testing.T) {
		t.Run("not found error", func(t *testing.T) {
			authorGetter := func(id values.CommentId) (core_values.UserId, error) {
				return "", core_errors.ErrNotFound
			}
			err := service.NewCommentLikeToggler(authorGetter, nil)(comment, caller)
			AssertError(t, err, client_errors.NotFound)
		})
		t.Run("some other error", func(t *testing.T) {
			authorGetter := func(values.CommentId) (core_values.UserId, error) {
				return "", RandomError()
			}
			err := service.NewCommentLikeToggler(authorGetter, nil)(comment, caller)
			AssertSomeError(t, err)
		})
	})
	likeToggler := func(target string, authorId, callerId core_values.UserId) error {
		if target == comment && authorId == author && callerId == caller {
			return nil
		}
		panic("unexected args")
	}
	t.Run("error case - like toggler throws, should FORWARD the error (since it can be a client error)", func(t *testing.T) {
		wantErr := RandomError()
		likeToggler := func(target string, authorId, callerId core_values.UserId) error {
			return wantErr
		}
		err := service.NewCommentLikeToggler(authorGetter, likeToggler)(comment, caller)
		AssertError(t, err, wantErr)
	})
	err := service.NewCommentLikeToggler(authorGetter, likeToggler)(comment, caller)
	AssertNoError(t, err)
}

func TestPostCommentsGetter(t *testing.T) {
	post := RandomString()
	caller := RandomId()
	comments := []entities.Comment{RandomComment()}
	contextedComments := []entities.ContextedComment{RandomContextedComment()}

	commentsGetter := func(postId post_values.PostId) ([]entities.Comment, error) {
		if postId == post {
			return comments, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting the comments returns some error", func(t *testing.T) {
		commentsGetter := func(post_values.PostId) ([]entities.Comment, error) {
			return []entities.Comment{}, RandomError()
		}
		_, err := service.NewPostCommentsGetter(commentsGetter, nil)(post, caller)
		AssertSomeError(t, err)
	})
	contextAdder := func(commentList []entities.Comment, callerId core_values.UserId) ([]entities.ContextedComment, error) {
		if reflect.DeepEqual(commentList, comments) && callerId == caller {
			return contextedComments, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - context adder returns some error", func(t *testing.T) {
		contextAdder := func([]entities.Comment, core_values.UserId) ([]entities.ContextedComment, error) {
			return nil, RandomError()
		}
		_, err := service.NewPostCommentsGetter(commentsGetter, contextAdder)(post, caller)
		AssertSomeError(t, err)
	})
	gotComments, err := service.NewPostCommentsGetter(commentsGetter, contextAdder)(post, caller)
	AssertNoError(t, err)
	Assert(t, gotComments, contextedComments, "returned comments")
}
