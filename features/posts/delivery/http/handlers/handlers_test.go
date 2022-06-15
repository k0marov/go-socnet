package handlers_test

import (
	"bytes"
	"context"
	"core/client_errors"
	"core/core_values"
	helpers "core/http_test_helpers"
	"core/ref"
	. "core/test_helpers"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"posts/delivery/http/handlers"
	"posts/domain/entities"
	"posts/domain/values"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
)

func createRequestWithProfileId(profileId core_values.UserId) *http.Request {
	request := helpers.CreateRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("profile_id", profileId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return request
}

func TestGetListById(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		randomProfile := RandomString()
		randomPosts := handlers.PostsResponse{
			Posts: []entities.Post{RandomPost(), RandomPost(), RandomPost()},
		}
		getter := func(profileId core_values.UserId) ([]entities.Post, error) {
			if profileId == randomProfile {
				return randomPosts.Posts, nil
			}
			panic("unexpected args")
		}
		request := createRequestWithProfileId(randomProfile)
		response := httptest.NewRecorder()
		handlers.NewGetListByIdHandler(getter).ServeHTTP(response, request)
		AssertJSONData(t, response, randomPosts)
	})
	t.Run("error case - profile id is not provided", func(t *testing.T) {
		request := helpers.CreateRequest(nil)
		response := httptest.NewRecorder()
		handlers.NewGetListByIdHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(core_values.UserId) ([]entities.Post, error) {
			return []entities.Post{}, err
		}
		request := createRequestWithProfileId("42")
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
func TestCreateNew_Parsing(t *testing.T) {
	createMultipartBody := func(text string, images []string) (io.Reader, string) {
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)
		defer writer.Close()

		writer.WriteField("text", text)
		for i, image := range images {
			fw, _ := writer.CreateFormFile("image_"+strconv.Itoa(i+1), RandomString())
			fw.Write([]byte(image))
		}
		return body, writer.FormDataContentType()
	}
	convertImages := func(images []string) []core_values.FileData {
		files := []core_values.FileData{}
		for _, image := range images {
			imageBytes := []byte(image)
			ref, _ := ref.NewRef(&imageBytes)
			files = append(files, ref)
		}
		return files
	}

	cases := []struct {
		text       string
		imagesData []string
	}{
		{"", []string{"Cat Image", "Sky Image"}},
		{"One Image", []string{"Puppy Image"}},
		// {"Five Images", []string{"1", "2", "3", "4", "5"}},
		{"Zero Images", []string{}},
	}

	for _, testCase := range cases {
		t.Run(testCase.text, func(t *testing.T) {
			author := RandomAuthUser()
			expectedNewPost := values.NewPostData{
				Author: author.Id,
				Text:   testCase.text,
				Images: convertImages(testCase.imagesData),
			}
			called := false
			creator := func(newPost values.NewPostData) error {
				if reflect.DeepEqual(newPost, expectedNewPost) {
					called = true
					return nil
				}
				panic(fmt.Sprintf("enexpected args: newPost = %+v", newPost))
			}
			requestBody, contentType := createMultipartBody(testCase.text, testCase.imagesData)
			request := helpers.AddAuthDataToRequest(helpers.CreateRequest(requestBody), author)
			request.Header.Set("Content-Type", contentType)
			response := httptest.NewRecorder()
			handlers.NewCreateNewHandler(creator).ServeHTTP(response, request)

			AssertStatusCode(t, response, http.StatusOK)
			Assert(t, called, true, "service called")
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
