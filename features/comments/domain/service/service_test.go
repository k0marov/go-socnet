package service_test

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"testing"
	"time"
)

func TestCommentCreator(t *testing.T) {
	newComment := RandomNewComment()
	createdId := RandomString()
	author := RandomContextedProfile()
	createdComment := entities.ContextedComment{
		Id:        createdId,
		Author:    author,
		Text:      newComment.Text,
		CreatedAt: time.Now(),
		Likes:     0,
		IsLiked:   false,
		IsMine:    true,
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
	comments := []entities.Comment{RandomComment()}
	author := RandomContextedProfile()
	caller := RandomId()
	isLiked := RandomBool()

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
		_, err := service.NewPostCommentsGetter(commentsGetter, nil, nil)(post, caller)
		AssertSomeError(t, err)
	})
	profileGetter := func(profileId core_values.UserId, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
		if profileId == comments[0].Author && callerId == caller {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting profile throws", func(t *testing.T) {
		profileGetter := func(profileId, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
			return profile_entities.ContextedProfile{}, RandomError()
		}
		_, err := service.NewPostCommentsGetter(commentsGetter, profileGetter, nil)(post, caller)
		AssertSomeError(t, err)
	})
	t.Run("error case - like checker throws", func(t *testing.T) {
		likeChecker := func(commentId values.CommentId, callerId core_values.UserId) (bool, error) {
			return false, RandomError()
		}
		_, err := service.NewPostCommentsGetter(commentsGetter, profileGetter, likeChecker)(post, caller)
		AssertSomeError(t, err)
	})
	t.Run("happy cases", func(t *testing.T) {
		addContext := func(comment entities.Comment, author profile_entities.ContextedProfile, isLiked, isMine bool) entities.ContextedComment {
			return entities.ContextedComment{
				Id:        comment.Id,
				Author:    author,
				Text:      comment.Text,
				CreatedAt: comment.CreatedAt,
				Likes:     comment.Likes,
				IsLiked:   isLiked,
				IsMine:    isMine,
			}
		}
		t.Run("isMine = false", func(t *testing.T) {
			likeChecker := func(commentId values.CommentId, callerId core_values.UserId) (bool, error) {
				if commentId == comments[0].Id && callerId == caller {
					return isLiked, nil
				}
				panic("unexpected args")
			}
			sut := service.NewPostCommentsGetter(commentsGetter, profileGetter, likeChecker)
			wantProfiles := []entities.ContextedComment{addContext(comments[0], author, isLiked, false)}
			gotProfiles, err := sut(post, caller)
			AssertNoError(t, err)
			Assert(t, gotProfiles, wantProfiles, "the returned profiles")
		})
		t.Run("isMine = true", func(t *testing.T) {
			profileGetter := func(target, caller core_values.UserId) (profile_entities.ContextedProfile, error) {
				return author, nil
			}
			likeChecker := func(commentId values.CommentId, callerId core_values.UserId) (bool, error) {
				return isLiked, nil
			}
			sut := service.NewPostCommentsGetter(commentsGetter, profileGetter, likeChecker)
			wantProfiles := []entities.ContextedComment{addContext(comments[0], author, isLiked, true)}
			gotProfiles, err := sut(post, author.Id)
			AssertNoError(t, err)
			Assert(t, gotProfiles, wantProfiles, "the returned profiles")
		})
	})
}
