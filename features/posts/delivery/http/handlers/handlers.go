package handlers

import (
	"core/client_errors"
	helpers "core/http_helpers"
	"encoding/json"
	"net/http"
	"posts/domain/entities"
	"posts/domain/service"

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
		helpers.SetJsonHeader(w)
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
		json.NewEncoder(w).Encode(PostsResponse{posts})
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

func NewCreateNewHandler(createNew service.PostCreater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
