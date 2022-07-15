package handlers_test

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	helpers "github.com/k0marov/go-socnet/core/helpers/http_test_helpers"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"github.com/k0marov/go-socnet/features/feed/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/feed/delivery/http/responses"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createRequestWithCount(count string, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodOptions, "/handler-should-not-care?count="+count, body)
}

func TestFeedHandler(t *testing.T) {
	count := RandomString()
	posts := []string{RandomString(), RandomString()}
	caller := RandomAuthUser()

	helpers.BaseTest401(t, handlers.NewFeedHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		getter := func(countStr string, callerId core_values.UserId) ([]string, error) {
			if countStr == count && callerId == caller.Id {
				return posts, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(createRequestWithCount(count, nil), caller)
		handlers.NewFeedHandler(getter).ServeHTTP(response, request)
		AssertJSONData(t, response, responses.FeedResponse{Posts: posts})
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		getter := func(string, core_values.UserId) ([]string, error) {
			return []string{}, err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithCount(count, nil), caller)
		handlers.NewFeedHandler(getter).ServeHTTP(response, request)
	})
}
