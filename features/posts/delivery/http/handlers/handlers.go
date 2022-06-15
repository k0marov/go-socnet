package handlers

import (
	"core/client_errors"
	"core/core_values"
	helpers "core/http_helpers"
	"net/http"
	"posts/domain/entities"
	"posts/domain/service"
	"posts/domain/values"
	"strconv"

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

type PostsResponse struct {
	Posts []entities.Post
}

func NewGetListByIdHandler(getPosts service.PostsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profileId := chi.URLParam(r, "profile_id")
		if profileId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		posts, err := getPosts(profileId)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, PostsResponse{posts})
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

func parseImages(r *http.Request) []core_values.FileData {
	images := []core_values.FileData{}
	for i := 1; ; i++ {
		image, ok := helpers.ParseFile(r, "image_"+strconv.Itoa(i))
		if !ok {
			return images
		}
		images = append(images, image)
	}
}
