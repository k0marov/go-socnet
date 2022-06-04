package server_test

import (
	"bytes"
	"core/client_errors"
	core_entities "core/entities"
	. "core/test_helpers"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"profiles/delivery/server"
	"profiles/domain/entities"
	"profiles/domain/values"
	"testing"
)

func TestHTTPServer_PUT_Me(t *testing.T) {
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
	goodAvatarPath := filepath.Join("testdata", "test_avatar.png")

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
				updateAvatar: func(u core_entities.User, avatar values.AvatarData) (entities.DetailedProfile, error) {
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

			AssertHTTPError(t, response, client_errors.AvatarNotProvidedError)
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

			AssertHTTPError(t, response, client_errors.BodyIsNotMultipartForm)
		})
	})
	baseTestServerErrorHandling(t, func() *http.Request { return createRequestWithAuth(goodAvatarPath, RandomString()) })
}
