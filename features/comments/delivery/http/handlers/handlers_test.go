package handlers_test

import (
	"context"
	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_values"
	helpers "github.com/k0marov/socnet/core/http_test_helpers"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/values"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGetCommentsHandler(t *testing.T) {
	createRequestWithPostId := func(id post_values.PostId) *http.Request {
		return httptest.NewRequest(http.MethodOptions, "/handler-should-not-care?post_id="+id, nil)
	}
	post := RandomString()
	comments := RandomComments()
	t.Run("happy case", func(t *testing.T) {
		getter := func(postId post_values.PostId) ([]entities.Comment, error) {
			if postId == post {
				return comments, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		handlers.NewGetCommentsHandler(getter).ServeHTTP(response, createRequestWithPostId(post))
		AssertJSONData(t, response, handlers.CommentsResponse{Comments: comments})
	})
	t.Run("error case - post_id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		handlers.NewGetCommentsHandler(nil).ServeHTTP(response, helpers.CreateRequest(nil))
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		getter := func(id post_values.PostId) ([]entities.Comment, error) {
			return []entities.Comment{}, err
		}
		handlers.NewGetCommentsHandler(getter).ServeHTTP(response, createRequestWithPostId(post))
	})
}

func TestNewAddCommentHandler(t *testing.T) {
}

func TestNewToggleLikeCommentHandler(t *testing.T) {
	createRequestWithCommentId := func(id values.CommentId, caller auth.User) *http.Request {
		request := helpers.CreateRequest(nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		return helpers.AddAuthDataToRequest(request, caller)
	}
	helpers.BaseTest401(t, handlers.NewToggleLikeCommentHandler(nil))

	user := RandomAuthUser()
	comment := RandomString()
	t.Run("happy case", func(t *testing.T) {
		likeToggler := func(commentId values.CommentId, caller core_values.UserId) error {
			if commentId == comment && caller == user.Id {
				return nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		handlers.NewToggleLikeCommentHandler(likeToggler).ServeHTTP(response, createRequestWithCommentId(comment, user))
		AssertStatusCode(t, response, http.StatusOK)
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		handlers.NewToggleLikeCommentHandler(nil).ServeHTTP(response, helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), user))
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		likeToggler := func(values.CommentId, core_values.UserId) error {
			return err
		}
		handlers.NewToggleLikeCommentHandler(likeToggler).ServeHTTP(response, createRequestWithCommentId(comment, user))
	})
}
