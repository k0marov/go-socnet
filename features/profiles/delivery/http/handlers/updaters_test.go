package handlers_test

import (
	"bytes"
	"core/client_errors"
	core_entities "core/entities"
	helpers "core/http_test_helpers"
	. "core/test_helpers"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/http/handlers"
	"profiles/domain/entities"
	"profiles/domain/values"
	"reflect"
	"testing"
)

func TestUpdateMeHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	profileUpdate := values.ProfileUpdateData{About: RandomString()}
	createGoodRequest := func() *http.Request {
		body := bytes.NewBuffer(nil)
		json.NewEncoder(body).Encode(profileUpdate)
		return helpers.AddAuthDataToRequest(helpers.CreateRequest(body), authUser)
	}

	helpers.BaseTest401(t, handlers.NewUpdateMeHandler(nil))
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
	helpers.BaseTestServiceErrorHandling(t, func(wantErr error, w *httptest.ResponseRecorder) {
		update := func(gotUser core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, wantErr
		}
		handlers.NewUpdateMeHandler(update).ServeHTTP(w, createGoodRequest())
	})
	t.Run("should return invalid json client error if request is not valid json", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(bytes.NewBufferString("non-json")), authUser)
		handler := handlers.NewUpdateMeHandler(nil) // service is nil, since it shouldn't be called
		handler.ServeHTTP(response, request)

		AssertClientError(t, response, client_errors.InvalidJsonError)
	})
}

func TestToggleFollowHandler(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewToggleFollowHandler(nil))
	t.Run("should toggle follow using service", func(t *testing.T) {
		targetId := RandomString()
		followerAuth := RandomAuthUser()
		called := false
		followToggler := func(target, follower values.UserId) error {
			if follower == followerAuth.Id && target == targetId {
				called = true
				return nil
			}
			panic("called with unexpected args")
		}

		request := helpers.AddAuthDataToRequest(createRequestWithId(targetId), followerAuth)
		response := httptest.NewRecorder()
		handlers.NewToggleFollowHandler(followToggler).ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		Assert(t, called, true, "toggle called")
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		handlers.NewToggleFollowHandler(nil).ServeHTTP(response, helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())) // toggler is nil, since it shouldn't be called
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		followToggler := func(target, follower values.UserId) error {
			return err
		}
		handlers.NewToggleFollowHandler(followToggler).ServeHTTP(w, helpers.AddAuthDataToRequest(createRequestWithId("42"), RandomAuthUser()))

	})
}

func TestUpdateAvatarHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	tAvatar := []byte(RandomString())

	getMultipartBody := func(data []byte) (io.Reader, string) {
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)
		defer writer.Close()
		fw, _ := writer.CreateFormFile("avatar", RandomString())
		fw.Write(data)
		return body, writer.FormDataContentType()
	}
	createRequestWithAuth := func() *http.Request {
		body, contentType := getMultipartBody(tAvatar)
		req := helpers.AddAuthDataToRequest(helpers.CreateRequest(body), authUser)
		req.Header.Set("Content-Type", contentType)
		return helpers.AddAuthDataToRequest(req, authUser)
	}
	helpers.BaseTest401(t, handlers.NewUpdateAvatarHandler(nil))
	t.Run("should update avatar using service", func(t *testing.T) {
		t.Run("happy case", func(t *testing.T) {
			avatarURL := values.AvatarPath{Path: RandomString()}
			updateAvatar := func(u core_entities.User, avatar values.AvatarData) (values.AvatarPath, error) {
				if u == user && reflect.DeepEqual(avatar.Data.Value(), tAvatar) {
					return avatarURL, nil
				}
				panic("updateAvatar called with improper arguments")
			}

			response := httptest.NewRecorder()
			handlers.NewUpdateAvatarHandler(updateAvatar).ServeHTTP(response, createRequestWithAuth())

			AssertStatusCode(t, response, http.StatusOK)
			AssertJSONData(t, response, avatarURL)
		})
		t.Run("error case - avatar file is not provided", func(t *testing.T) {
			response := httptest.NewRecorder()
			req := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), authUser)
			handler := handlers.NewUpdateAvatarHandler(nil) // since the service function shouldn't be called, it's nil
			handler.ServeHTTP(response, req)

			AssertClientError(t, response, client_errors.AvatarNotProvidedError)
		})
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		updateAvatar := func(core_entities.User, values.AvatarData) (values.AvatarPath, error) {
			return values.AvatarPath{}, err
		}
		handlers.NewUpdateAvatarHandler(updateAvatar).ServeHTTP(w, createRequestWithAuth())
	})
}
