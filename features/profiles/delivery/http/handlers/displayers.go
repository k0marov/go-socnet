package handlers

import (
	"core/client_errors"
	"core/core_values"
	helpers "core/http_helpers"
	"net/http"
	"profiles/domain/service"

	"github.com/go-chi/chi/v5"
)

func NewGetMeHandler(getProfile service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		profile, err := getProfile(user.Id, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, profile)
	})
}

func NewGetByIdHandler(profileGetter service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		profile, err := profileGetter(id, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, profile)
	})
}

type FollowsResponse struct {
	Profiles []core_values.UserId
}

func NewGetFollowsHandler(followsGetter service.FollowsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		follows, err := followsGetter(id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, FollowsResponse{follows})
	})
}
