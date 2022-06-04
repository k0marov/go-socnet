package server_test

import (
	"context"
	"core/client_errors"
	core_entities "core/entities"
	. "core/test_helpers"
	"errors"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/server"
	"profiles/domain/entities"
	"profiles/domain/values"
	"testing"

	auth "github.com/k0marov/golang-auth"
)

var dummyProfileService = &StubProfileService{}

func addAuthDataToRequest(r *http.Request, user auth.User) *http.Request {
	ctx := context.WithValue(r.Context(), auth.UserContextKey, user)
	return r.WithContext(ctx)
}

func baseTest401(t *testing.T, createRequest func() *http.Request) {
	t.Helper()
	t.Run("should return 401 if authentication details are not provided via context (using auth middleware)", func(t *testing.T) {
		srv := server.NewHTTPServer(dummyProfileService)

		request := createRequest()
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusUnauthorized)
	})
}

func baseTestServerErrorHandling(t *testing.T, createRequestWithAuth func() *http.Request) {
	t.Helper()

	t.Run("service throws an error", func(t *testing.T) {
		t.Run("it is a client error - encode and deliver it to client", func(t *testing.T) {
			randomClientError := client_errors.ClientError{
				DetailCode:     RandomString(),
				ReadableDetail: RandomString(),
				HTTPCode:       RandomInt() + 400,
			}
			srv := server.NewHTTPServer(&StubProfileService{returnedError: randomClientError, doNotPanic: true})

			request := createRequestWithAuth()
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			AssertHTTPError(t, response, randomClientError)
		})
		t.Run("it is a server error - just return status code 500", func(t *testing.T) {
			randomError := errors.New(RandomString())
			srv := server.NewHTTPServer(&StubProfileService{returnedError: randomError, doNotPanic: true})

			request := createRequestWithAuth()
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			AssertStatusCode(t, response, http.StatusInternalServerError)
		})
	})
}

func TestHTTPServer_Me_IncorrectMethod(t *testing.T) {
	srv := server.NewHTTPServer(dummyProfileService)
	response := httptest.NewRecorder()
	srv.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/profiles/me", nil))

	AssertStatusCode(t, response, http.StatusMethodNotAllowed)
}

type StubProfileService struct {
	getDetailed   func(core_entities.User) (entities.DetailedProfile, error)
	update        func(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
	updateAvatar  func(user core_entities.User, avatar values.AvatarData) (entities.DetailedProfile, error)
	returnedError error
	doNotPanic    bool
}

func (s *StubProfileService) GetDetailed(user core_entities.User) (entities.DetailedProfile, error) {
	if s.getDetailed != nil {
		return s.getDetailed(user)
	}
	if s.doNotPanic {
		return entities.DetailedProfile{}, s.returnedError
	}
	panic("getOrCreate method shouldn't have been called")
}
func (s *StubProfileService) Update(user core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
	if s.update != nil {
		return s.update(user, updateData)
	}
	if s.doNotPanic {
		return entities.DetailedProfile{}, s.returnedError
	}
	panic("Update method shouldn't have been called")
}
func (s *StubProfileService) UpdateAvatar(user core_entities.User, avatar values.AvatarData) (entities.DetailedProfile, error) {
	if s.updateAvatar != nil {
		return s.updateAvatar(user, avatar)
	}
	if s.doNotPanic {
		return entities.DetailedProfile{}, s.returnedError
	}
	panic("UpdateAvatar method shouldn't have been called")
}
