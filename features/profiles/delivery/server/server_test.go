package server_test

import (
	"bytes"
	"context"
	"core/client_errors"
	core_entities "core/entities"
	. "core/test_helpers"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
			}
			srv := server.NewHTTPServer(&StubProfileService{returnedError: randomClientError, doNotPanic: true})

			request := createRequestWithAuth()
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			AssertHTTPError(t, response, randomClientError, http.StatusBadRequest)
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

func TestHTTPServer_GET_Me(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	createRequest := func() *http.Request {
		return httptest.NewRequest(http.MethodGet, "/profiles/me", nil)
	}
	createRequestWithAuth := func() *http.Request {
		r := createRequest()
		return addAuthDataToRequest(r, authUser)
	}

	baseTest401(t, createRequest)
	t.Run("should return 200 and a profile if authentication details are provided via context", func(t *testing.T) {
		wantedProfile := RandomDetailedProfile()
		service := &StubProfileService{
			getOrCreateDetailed: func(u core_entities.User) (entities.DetailedProfile, error) {
				if u == user {
					return wantedProfile, nil
				}
				panic("getOrCreate called with improper arguments")
			},
		}
		srv := server.NewHTTPServer(service)

		request := createRequestWithAuth()
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, wantedProfile)
	})
	baseTestServerErrorHandling(t, createRequestWithAuth)
}

func TestHTTPServer_Me_IncorrectMethod(t *testing.T) {
	srv := server.NewHTTPServer(dummyProfileService)
	response := httptest.NewRecorder()
	srv.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/profiles/me", nil))

	AssertStatusCode(t, response, http.StatusMethodNotAllowed)
}

func TestHTTPServer_POST_Me(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	profileUpdate := values.ProfileUpdateData{
		About: RandomString(),
	}
	createRequest := func() *http.Request {
		body := bytes.NewBuffer(nil)
		json.NewEncoder(body).Encode(profileUpdate)
		return httptest.NewRequest(http.MethodPut, "/profiles/me", body)
	}
	createRequestWithAuth := func() *http.Request {
		r := createRequest()
		return addAuthDataToRequest(r, authUser)
	}
	baseTest401(t, createRequest)
	t.Run("should update profile about if about field is provided", func(t *testing.T) {
		updatedProfile := RandomDetailedProfile()
		service := &StubProfileService{
			update: func(u core_entities.User, pud values.ProfileUpdateData) (entities.DetailedProfile, error) {
				if u == user && pud == profileUpdate {
					return updatedProfile, nil
				}
				panic("update called with improper arguments")
			},
		}
		srv := server.NewHTTPServer(service)

		response := httptest.NewRecorder()
		srv.ServeHTTP(response, createRequestWithAuth())

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, updatedProfile)
	})
	baseTestServerErrorHandling(t, createRequestWithAuth)
}

func TestHTTPServer_Put_Me_Avatar(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	goodAvatarPath := "testdata/test_avatar.png"
	bigAvatarPath := "testdata/test_big_avatar.png"

	createRequest := func(avatarFilePath string, fileName string) *http.Request {
		body := bytes.NewBuffer(nil)

		writer := multipart.NewWriter(body)
		if avatarFilePath != "" {
			fw, _ := writer.CreateFormFile("avatar", fileName)
			file, err := os.Open(avatarFilePath)
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(fw, file)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/profiles/me/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		return req
	}
	createRequestWithAuth := func(avatarFilePath, fileName string) *http.Request {
		r := createRequest(avatarFilePath, fileName)
		return addAuthDataToRequest(r, authUser)
	}

	baseTest401(t, func() *http.Request { return createRequest(goodAvatarPath, RandomString()) })

	t.Run("should update avatar using service", func(t *testing.T) {
		t.Run("happy case", func(t *testing.T) {
			updatedProfile := RandomDetailedProfile()
			fileName := RandomString()
			service := &StubProfileService{
				updateAvatar: func(u core_entities.User, avatar server.AvatarData) (entities.DetailedProfile, error) {
					if u == user && avatar.FileName == fileName {
						return updatedProfile, nil
					}
					panic("updateAvatar called with improper arguments")
				},
			}
			srv := server.NewHTTPServer(service)

			response := httptest.NewRecorder()
			srv.ServeHTTP(response, createRequestWithAuth(goodAvatarPath, fileName))

			AssertStatusCode(t, response, http.StatusOK)
			AssertJSONData(t, response, updatedProfile)
		})
		t.Run("error case - avatar file is not provided", func(t *testing.T) {
			service := &StubProfileService{} // should not be called, so it panics on every call
			srv := server.NewHTTPServer(service)

			response := httptest.NewRecorder()
			req := createRequestWithAuth("", "")
			srv.ServeHTTP(response, req)

			AssertHTTPError(t, response, client_errors.AvatarNotProvidedError, http.StatusBadRequest)
		})
		t.Run("error case - avatar is too big", func(t *testing.T) {
			service := &StubProfileService{}
			srv := server.NewHTTPServer(service)

			response := httptest.NewRecorder()
			req := createRequestWithAuth(bigAvatarPath, RandomString())
			srv.ServeHTTP(response, req)

			AssertHTTPError(t, response, client_errors.AvatarTooBigError, http.StatusBadRequest)
		})
		t.Run("error case - post body is not multipart form", func(t *testing.T) {
			service := &StubProfileService{}
			srv := server.NewHTTPServer(service)

			someRandomData := values.ProfileUpdateData{About: RandomString()}
			dataJson, _ := json.Marshal(someRandomData)
			body := bytes.NewBuffer(dataJson)
			response := httptest.NewRecorder()
			request := addAuthDataToRequest(httptest.NewRequest(http.MethodPut, "/profiles/me/avatar", body), authUser)

			srv.ServeHTTP(response, request)

			AssertHTTPError(t, response, client_errors.BodyIsNotMultipartForm, http.StatusBadRequest)
		})
	})
	baseTestServerErrorHandling(t, func() *http.Request { return createRequestWithAuth(goodAvatarPath, RandomString()) })
}

type StubProfileService struct {
	getOrCreateDetailed func(core_entities.User) (entities.DetailedProfile, error)
	update              func(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
	updateAvatar        func(user core_entities.User, avatar server.AvatarData) (entities.DetailedProfile, error)
	returnedError       error
	doNotPanic          bool
}

func (s *StubProfileService) GetOrCreateDetailed(user core_entities.User) (entities.DetailedProfile, error) {
	if s.getOrCreateDetailed != nil {
		return s.getOrCreateDetailed(user)
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
func (s *StubProfileService) UpdateAvatar(user core_entities.User, avatar server.AvatarData) (entities.DetailedProfile, error) {
	if s.updateAvatar != nil {
		return s.updateAvatar(user, avatar)
	}
	if s.doNotPanic {
		return entities.DetailedProfile{}, s.returnedError
	}
	panic("UpdateAvatar method shouldn't have been called")
}
