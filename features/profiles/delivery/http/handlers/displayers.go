package handlers

import (
	"core/client_errors"
	helpers "core/http_helpers"
	"encoding/json"
	"net/http"
	"profiles/domain/entities"
	"profiles/domain/service"

	"github.com/go-chi/chi/v5"
)

func NewGetMeHandler(detailedProfileGetter service.DetailedProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		helpers.SetJsonHeader(w)
		profile, err := detailedProfileGetter(user)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		json.NewEncoder(w).Encode(profile)
	})
}

func NewGetByIdHandler(profileGetter service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.SetJsonHeader(w)
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		profile, err := profileGetter(id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		json.NewEncoder(w).Encode(profile)
	})
}

type FollowsResponse struct {
	Profiles []entities.Profile
}

func NewGetFollowsHandler(followsGetter service.FollowsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.SetJsonHeader(w)
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
		followsResp := FollowsResponse{follows}
		json.NewEncoder(w).Encode(followsResp)
	})
}
