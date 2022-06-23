package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/k0marov/socnet/features/posts/delivery/http/handlers"
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"

	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_values"
	helpers "github.com/k0marov/socnet/core/http_test_helpers"
	. "github.com/k0marov/socnet/core/test_helpers"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

func TestGetListById(t *testing.T) {
	caller := RandomAuthUser()
	helpers.BaseTest401(t, handlers.NewGetListByIdHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		randomProfile := RandomString()
		posts := []entities.ContextedPost{RandomContextedPost(), RandomContextedPost()}
		getter := func(profileId, callerId core_values.UserId) ([]entities.ContextedPost, error) {
			if profileId == randomProfile && callerId == caller.Id {
				return posts, nil
			}
			panic("unexpected args")
		}
		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodGet, "/handler-should-not-care?profile_id="+randomProfile, nil), caller)
		response := httptest.NewRecorder()
		handlers.NewGetListByIdHandler(getter).ServeHTTP(response, request)
		AssertJSONData(t, response, handlers.PostsToResponse(posts))
	})
	t.Run("error case - profile id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), caller)
		response := httptest.NewRecorder()
		handlers.NewGetListByIdHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(profile, caller core_values.UserId) ([]entities.ContextedPost, error) {
			return []entities.ContextedPost{}, err
		}
		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodGet, "/handler-should-not-care?profile_id=42", nil), caller)
		handlers.NewGetListByIdHandler(getter).ServeHTTP(rr, request)
	})
}
func TestToggleLike(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewToggleLikeHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		randomPost := RandomString()
		randomUser := RandomAuthUser()
		called := false
		toggler := func(post values.PostId, fromUser core_values.UserId) error {
			if post == randomPost && fromUser == randomUser.Id {
				called = true
				return nil
			}
			panic("unexpected args")
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(randomPost), randomUser)
		response := httptest.NewRecorder()
		handlers.NewToggleLikeHandler(toggler).ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK)
		Assert(t, called, true, "service called")
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		response := httptest.NewRecorder()
		handlers.NewToggleLikeHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		toggler := func(values.PostId, core_values.UserId) error {
			return err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId("42"), RandomAuthUser())
		handlers.NewToggleLikeHandler(toggler).ServeHTTP(rr, request)
	})
}

func TestCreatePost_ErrorHandling(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewCreateHandler(nil))
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		creator := func(values.NewPostData) error {
			return err
		}
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		handlers.NewCreateHandler(creator).ServeHTTP(rr, request)
	})
}

func TestCreatePost_Parsing(t *testing.T) {
	createRequest := func(postData values.NewPostData) *http.Request {
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)

		writer.WriteField("text", postData.Text)
		for _, image := range postData.Images {
			fw, _ := writer.CreateFormFile(fmt.Sprintf("image_%d", image.Index), RandomString())
			fw.Write(image.File.Value())
		}
		writer.Close()

		user := auth.User{Id: postData.Author, Username: RandomString()}
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(body), user)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		return request
	}

	cases := []values.NewPostData{
		{
			Text:   "0 Images",
			Author: "42",
			Images: []values.PostImageFile{},
		},
		{
			Text:   "2 Images",
			Author: "77",
			Images: []values.PostImageFile{{RandomFileData(), 1}, {RandomFileData(), 2}},
		},
		{
			Text:   "3 images",
			Author: "33",
			Images: []values.PostImageFile{{RandomFileData(), 1}, {RandomFileData(), 2}, {RandomFileData(), 3}},
		},
	}

	for _, testNewPost := range cases {
		t.Run(testNewPost.Text, func(t *testing.T) {
			creator := func(newPost values.NewPostData) error {
				if reflect.DeepEqual(newPost, testNewPost) {
					return nil
				}
				panic(fmt.Sprintf("enexpected args: newPost = %+v", newPost))
			}
			response := httptest.NewRecorder()
			handlers.NewCreateHandler(creator).ServeHTTP(response, createRequest(testNewPost))

			AssertStatusCode(t, response, http.StatusOK)
		})
	}
}

func createRequestWithPostId(postId values.PostId) *http.Request {
	request := helpers.CreateRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", postId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return request
}
func TestDelete(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewDeleteHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		randomPost := RandomString()
		randomUser := RandomAuthUser()
		called := false
		deleter := func(post values.PostId, fromUser core_values.UserId) error {
			if post == randomPost && fromUser == randomUser.Id {
				called = true
				return nil
			}
			panic("unexpected args")
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId(randomPost), randomUser)
		response := httptest.NewRecorder()
		handlers.NewDeleteHandler(deleter).ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		Assert(t, called, true, "service was called")
	})
	t.Run("error case - post id is not provided", func(t *testing.T) {
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		response := httptest.NewRecorder()
		handlers.NewDeleteHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		deleter := func(values.PostId, core_values.UserId) error {
			return err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithPostId("42"), RandomAuthUser())
		handlers.NewDeleteHandler(deleter).ServeHTTP(response, request)
	})
}
