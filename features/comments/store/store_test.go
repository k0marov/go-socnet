package store_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	comment_models "github.com/k0marov/socnet/features/comments/domain/models"
	"github.com/k0marov/socnet/features/comments/store"
	"testing"
)

func TestCommentsGetter(t *testing.T) {
	commentModels := []comment_models.CommentModel{RandomCommentModel()}
	likes := RandomInt()
	author := RandomId()

	commentsGetter := func(authorId core_values.UserId) ([]comment_models.CommentModel, error) {
		if authorId == author {
			return commentModels, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting comments from db throws", func(t *testing.T) {
		commentsGetter := func(id core_values.UserId) ([]comment_models.CommentModel, error) {
			return nil, RandomError()
		}
		_, err := store.NewCommentsGetter(commentsGetter, nil)(author)
		AssertSomeError(t, err)
	})
	likesGetter := func(targetId string) (int, error) {
		if targetId == commentModels[0].Id {
			return likes, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting likes throws", func(t *testing.T) {
		likesGetter := func(string) (int, error) {
			return 0, RandomError()
		}
		_, err := store.NewCommentsGetter(commentsGetter, likesGetter)(author)
		AssertSomeError(t, err)
	})
	gotComments, err := store.NewCommentsGetter(commentsGetter, likesGetter)(author)
	AssertNoError(t, err)
	wantComments := []entities.Comment{
		{
			CommentModel: commentModels[0],
			Likes:        likes,
		},
	}
	Assert(t, gotComments, wantComments, "returned comments")
}
