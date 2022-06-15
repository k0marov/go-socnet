package handlers_test

import (
	"context"
	"core/client_errors"
	"core/core_values"
	helpers "core/http_test_helpers"
	. "core/test_helpers"
	"net/http"
	"net/http/httptest"
	"posts/delivery/http/handlers"
	"posts/domain/values"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetListById(t *testing.T) {
}
func TestToggleLike(t *testing.T) {

}
func TestCreateNew(t *testing.T) {

}

func createRequestWithPostId(postId values.PostId) *http.Request {
	request := helpers.CreateRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", postId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return request
}
func TestDelete(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewDeleteHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		randomPost := RandomString()
		randomUser := RandomAuthUser()
		called := false
		deleter := func(post values.PostId, fromUser core_values.UserId) error {
			if post == randomPost && fromUser == randomUser.Id {
				called = true
				return nil
			}
			panic("unexpected args")
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(randomPost), randomUser)
		response := httptest.NewRecorder()
		handlers.NewDeleteHandler(deleter).ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		Assert(t, called, true, "service was called")
	})
	t.Run("error case - post id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		response := httptest.NewRecorder()
		handlers.NewDeleteHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		deleter := func(values.PostId, core_values.UserId) error {
			return err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId("42"), RandomAuthUser())
		handlers.NewDeleteHandler(deleter).ServeHTTP(response, request)
	})
}
