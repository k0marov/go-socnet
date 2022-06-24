package handlers

import (
	"net/http"
	"strconv"

	"github.com/k0marov/go-socnet/features/posts/delivery/http/responses"
	"github.com/k0marov/go-socnet/features/posts/domain/service"
	"github.com/k0marov/go-socnet/features/posts/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	helpers "github.com/k0marov/go-socnet/core/http_helpers"

	"github.com/go-chi/chi/v5"
)

func NewDeleteHandler(deletePost service.PostDeleter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := chi.URLParam(r, "id")
		if postId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		err := deletePost(postId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
	})
}

func NewGetListByIdHandler(getPosts service.PostsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		profileId := r.URL.Query().Get("profile_id")
		if profileId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		posts, err := getPosts(profileId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, responses.NewPostListResponse(posts))
	})
}

func NewToggleLikeHandler(toggleLike service.PostLikeToggler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := chi.URLParam(r, "id")
		if postId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		err := toggleLike(postId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
		}
	})
}

func NewCreateHandler(createPost service.PostCreator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		newPost := values.NewPostData{
			Author: user.Id,
			Text:   r.FormValue("text"),
			Images: parseImages(r),
		}
		err := createPost(newPost)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
	})
}

func parseImages(r *http.Request) []values.PostImageFile {
	images := []values.PostImageFile{}
	for i := 1; ; i++ {
		file, ok := helpers.ParseFile(r, "image_"+strconv.Itoa(i))
		if !ok {
			return images
		}
		image := values.PostImageFile{File: file, Index: i}
		images = append(images, image)
	}
}
