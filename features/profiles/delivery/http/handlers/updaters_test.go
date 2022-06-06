package handlers_test

import (
	"bytes"
	"core/client_errors"
	core_entities "core/entities"
	. "core/test_helpers"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"profiles/delivery/http/handlers"
	"profiles/domain/entities"
	"profiles/domain/values"
	"testing"
)

func TestUpdateMeHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	profileUpdate := values.ProfileUpdateData{About: RandomString()}
	createGoodRequest := func() *http.Request {
		body := bytes.NewBuffer(nil)
		json.NewEncoder(body).Encode(profileUpdate)
		return addAuthDataToRequest(createRequest(body), authUser)
	}

	baseTest401(t, handlers.NewUpdateMeHandler(nil))
	t.Run("should update profile about if about field is provided", func(t *testing.T) {
		updatedProfile := RandomDetailedProfile()
		update := func(gotUser core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
			if gotUser == user && updateData == profileUpdate {
				return updatedProfile, nil
			}
			panic(fmt.Sprintf("called with gotUser=%v and updateData=%v", gotUser, updateData))
		}

		response := httptest.NewRecorder()
		handlers.NewUpdateMeHandler(update).ServeHTTP(response, createGoodRequest())

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, updatedProfile)
	})
	baseTestServiceErrorHandling(t, func(wantErr error, w *httptest.ResponseRecorder) {
		update := func(gotUser core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, wantErr
		}
		handlers.NewUpdateMeHandler(update).ServeHTTP(w, createGoodRequest())
	})
	t.Run("should return invalid json client error if request is not valid json", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := addAuthDataToRequest(createRequest(bytes.NewBufferString("non-json")), authUser)
		handler := handlers.NewUpdateMeHandler(nil) // service is nil, since it shouldn't be called
		handler.ServeHTTP(response, request)

		AssertClientError(t, response, client_errors.InvalidJsonError)
	})
}

func TestUpdateAvatarHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	goodAvatarPath := filepath.Join("testdata", "test_avatar.jpg")

	getMultipartBody := func(filePath, fileName string) (io.Reader, string) {
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)
		defer writer.Close()
		if filePath != "" {
			fw, _ := writer.CreateFormFile("avatar", fileName)
			file, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("Error while opening fixture file: %v", err)
			}
			io.Copy(fw, file)
		}
		return body, writer.FormDataContentType()
	}
	createRequestWithAuth := func(avatarFilePath, fileName string) *http.Request {
		body, contentType := getMultipartBody(avatarFilePath, fileName)
		req := addAuthDataToRequest(createRequest(body), authUser)
		req.Header.Set("Content-Type", contentType)
		return addAuthDataToRequest(req, authUser)
	}
	baseTest401(t, handlers.NewUpdateAvatarHandler(nil))
	t.Run("should update avatar using service", func(t *testing.T) {
		t.Run("happy case", func(t *testing.T) {
			fileName := RandomString()
			avatarURL := values.AvatarURL{Url: RandomString()}
			updateAvatar := func(u core_entities.User, avatar values.AvatarData) (values.AvatarURL, error) {
				if u == user && avatar.FileName == fileName {
					return avatarURL, nil
				}
				panic("updateAvatar called with improper arguments")
			}

			response := httptest.NewRecorder()
			handlers.NewUpdateAvatarHandler(updateAvatar).ServeHTTP(response, createRequestWithAuth(goodAvatarPath, fileName))

			AssertStatusCode(t, response, http.StatusOK)
			AssertJSONData(t, response, avatarURL)
		})
		t.Run("error case - avatar file is not provided", func(t *testing.T) {
			response := httptest.NewRecorder()
			req := addAuthDataToRequest(createRequest(nil), authUser)
			handler := handlers.NewUpdateAvatarHandler(nil) // since the service function shouldn't be called, it's nil
			handler.ServeHTTP(response, req)

			AssertClientError(t, response, client_errors.AvatarNotProvidedError)
		})
	})
	baseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		updateAvatar := func(core_entities.User, values.AvatarData) (values.AvatarURL, error) {
			return values.AvatarURL{}, err
		}
		handlers.NewUpdateAvatarHandler(updateAvatar).ServeHTTP(w, createRequestWithAuth(goodAvatarPath, RandomString()))
	})
}
