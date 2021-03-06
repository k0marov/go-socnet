package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	helpers "github.com/k0marov/go-socnet/core/helpers/http_test_helpers"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/responses"
	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
	auth "github.com/k0marov/golang-auth"
)

func createRequestWithPostId(id post_values.PostId, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodOptions, "/handler-should-not-care?post_id="+id, body)
}

func TestNewGetCommentsHandler(t *testing.T) {
	post := RandomString()
	comments := []entities.ContextedComment{RandomContextedComment()}
	caller := RandomAuthUser()

	helpers.BaseTest401(t, handlers.NewGetCommentsHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		getter := func(postId post_values.PostId, callerId core_values.UserId) ([]entities.ContextedComment, error) {
			if postId == post && callerId == caller.Id {
				return comments, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(post, nil), caller)
		handlers.NewGetCommentsHandler(getter).ServeHTTP(response, request)
		AssertJSONData(t, response, responses.NewCommentListResponse(comments))
	})
	t.Run("error case - post_id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), caller)
		handlers.NewGetCommentsHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		getter := func(post_values.PostId, core_values.UserId) ([]entities.ContextedComment, error) {
			return []entities.ContextedComment{}, err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(post, nil), caller)
		handlers.NewGetCommentsHandler(getter).ServeHTTP(response, request)
	})
}

func TestNewCreateCommentHandler(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewCreateCommentHandler(nil))
	user := RandomAuthUser()
	post := RandomString()
	wantNewComment := values.NewCommentValue{
		Author: user.Id,
		Text:   RandomString(),
		Post:   post,
	}
	createdComment := RandomContextedComment()
	t.Run("happy case", func(t *testing.T) {
		creator := func(newComment values.NewCommentValue) (entities.ContextedComment, error) {
			if newComment == wantNewComment {
				return createdComment, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		body := bytes.NewBuffer(nil)
		json.NewEncoder(body).Encode(handlers.NewCommentRequest{Text: wantNewComment.Text})
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(post, body), user)
		handlers.NewCreateCommentHandler(creator).ServeHTTP(response, request)
		AssertJSONData(t, response, responses.NewCommentResponse(createdComment))
	})
	t.Run("error case - post id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), user)
		response := httptest.NewRecorder()
		handlers.NewCreateCommentHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		creator := func(values.NewCommentValue) (entities.ContextedComment, error) {
			return entities.ContextedComment{}, err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(post, nil), user)
		handlers.NewCreateCommentHandler(creator).ServeHTTP(response, request)
	})
}

func createRequestWithCommentId(id values.CommentId, caller auth.User) *http.Request {
	request := helpers.CreateRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", id)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return helpers.AddAuthDataToRequest(request, caller)
}

func TestNewToggleLikeCommentHandler(t *testing.T) {
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

func TestDeleteCommentHandler(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewDeleteCommentHandler(nil))
	comment := RandomId()
	caller := RandomAuthUser()
	t.Run("happy case", func(t *testing.T) {
		delete := func(commentId values.CommentId, callerId core_values.UserId) error {
			if commentId == comment && callerId == caller.Id {
				return nil
			}
			panic("unexpected args")
		}
		request := createRequestWithCommentId(comment, caller)
		response := httptest.NewRecorder()
		handlers.NewDeleteCommentHandler(delete)(response, request)

		AssertStatusCode(t, response, http.StatusOK)
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), caller)
		response := httptest.NewRecorder()
		handlers.NewDeleteCommentHandler(nil)(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		delete := func(values.CommentId, core_values.UserId) error {
			return err
		}
		request := createRequestWithCommentId(comment, caller)
		handlers.NewDeleteCommentHandler(delete)(response, request)
	})
}
