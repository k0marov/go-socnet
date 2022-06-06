package handlers_test

import (
	"context"
	. "core/test_helpers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/k0marov/golang-auth"
)

func addAuthDataToRequest(r *http.Request, user auth.User) *http.Request {
	ctx := context.WithValue(r.Context(), auth.UserContextKey, user)
	return r.WithContext(ctx)
}

// Since I am using go-chi router, handlers can be independent of urls and http methods
func createRequest(body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/handler-should-not-care", body)
}

func baseTestServiceErrorHandling(t *testing.T, callErroringHandler func(error, *httptest.ResponseRecorder)) {
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

func baseTest401(t *testing.T, handlerWithPanickingService http.Handler) {
	t.Helper()
	t.Run("should return 401 if authentication details are not provided via context (using auth middleware)", func(t *testing.T) {
		request := createRequest(nil)
		response := httptest.NewRecorder()
		handlerWithPanickingService.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusUnauthorized)
	})
}
