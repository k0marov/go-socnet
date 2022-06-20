package handlers_test

import (
	helpers "github.com/k0marov/socnet/core/http_test_helpers"
	"net/http/httptest"
	"testing"
)

func TestNewGetCommentsHandler(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {

	})
	t.Run("error case - post_id is not provided", func(t *testing.T) {

	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {

	})
}

func TestNewAddCommentHandler(t *testing.T) {
}

func TestNewToggleLikeCommentHandler(t *testing.T) {

}
