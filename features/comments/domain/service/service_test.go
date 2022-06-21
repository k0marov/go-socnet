package service_test

import (
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/store/models"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"testing"
)

func TestCommentCreator(t *testing.T) {

}

func TestCommentLikeToggler(t *testing.T) {

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
