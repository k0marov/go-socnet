package http_test_helpers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/k0marov/go-socnet/core/test_helpers"

	auth "github.com/k0marov/golang-auth"
)

func AddAuthDataToRequest(r *http.Request, user auth.User) *http.Request {
	ctx := context.WithValue(r.Context(), auth.UserContextKey, user)
	return r.WithContext(ctx)
}

// Since I am using go-chi router, handlers can be independent of urls and http methods
func CreateRequest(body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/handler-should-not-care", body)
}

func BaseTestServiceErrorHandling(t *testing.T, callErroringHandler func(error, *httptest.ResponseRecorder)) {
	t.Helper()
	t.Run("service throws a client error", func(t *testing.T) {
		clientError := RandomClientError()
		response := httptest.NewRecorder()
		callErroringHandler(clientError, response)
		AssertClientError(t, response, clientError)
	})
	t.Run("service throws an internal error", func(t *testing.T) {
		response := httptest.NewRecorder()
		callErroringHandler(RandomError(), response)
		AssertStatusCode(t, response, http.StatusInternalServerError)
	})
}

func BaseTest401(t *testing.T, handlerWithPanickingService http.Handler) {
	t.Helper()
	t.Run("should return 401 if authentication details are not provided via context (using auth middleware)", func(t *testing.T) {
		request := CreateRequest(nil)
		response := httptest.NewRecorder()
		handlerWithPanickingService.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusUnauthorized)
	})
}
