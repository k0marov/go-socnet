package server_test

import (
	"context"
	"core/client_errors"
	. "core/test_helpers"
	"errors"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/server"
	"profiles/domain/entities"
	"testing"

	auth "github.com/k0marov/golang-auth"
)

type StubProfileService struct {
	returnedProfile entities.Profile
	returnedError   error
}

func (s *StubProfileService) GetOrCreate(user auth.User) (entities.Profile, error) {
	return s.returnedProfile, s.returnedError
}

var dummyProfileService = &StubProfileService{}

func TestHTTPServer_GET_Me(t *testing.T) {
	createRequest := func() *http.Request {
		return httptest.NewRequest(http.MethodGet, "/profiles/me", nil)
	}
	createRequestWithAuth := func() *http.Request {
		user := auth.User{
			Id:       "42",
			Username: "Sam",
		}
		r := createRequest()
		ctx := context.WithValue(r.Context(), "User", user)
		return r.WithContext(ctx)
	}

	t.Run("should return 401 if authentication details are not provided via context (using auth middleware)", func(t *testing.T) {
		srv := server.NewHTTPServer(dummyProfileService)

		request := createRequest()
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusUnauthorized)
	})
	t.Run("should return 200 and a profile if authentication details are provided via context", func(t *testing.T) {
		wantedProfile := entities.Profile{
			Id:       RandomString(),
			Username: RandomString(),
		}
		srv := server.NewHTTPServer(&StubProfileService{wantedProfile, nil})

		request := createRequestWithAuth()
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, wantedProfile)
	})
	t.Run("service throws an error", func(t *testing.T) {
		t.Run("it is a client error - encode and deliver it to client", func(t *testing.T) {
			randomClientError := client_errors.ClientError{
				DetailCode:     RandomString(),
				ReadableDetail: RandomString(),
			}
			srv := server.NewHTTPServer(&StubProfileService{returnedError: randomClientError})

			request := createRequestWithAuth()
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			AssertHTTPError(t, response, randomClientError, http.StatusBadRequest)
		})
		t.Run("it is a server error - just return status code 500", func(t *testing.T) {
			randomError := errors.New(RandomString())
			srv := server.NewHTTPServer(&StubProfileService{returnedError: randomError})

			request := createRequestWithAuth()
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			AssertStatusCode(t, response, http.StatusInternalServerError)
		})
	})
}
