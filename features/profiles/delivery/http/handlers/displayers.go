package handlers

import (
	"net/http"

	"github.com/k0marov/go-socnet/features/profiles/delivery/http/responses"

	"github.com/k0marov/go-socnet/features/profiles/domain/service"

	"github.com/k0marov/go-socnet/core/client_errors"
	helpers "github.com/k0marov/go-socnet/core/http_helpers"

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
		helpers.WriteJson(w, responses.NewProfileResponse(profile))
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
		helpers.WriteJson(w, responses.NewProfileResponse(profile))
	})
}

func NewGetFollowsHandler(followsGetter service.FollowsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		follows, err := followsGetter(id, caller.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, responses.NewProfilesResponse(follows))
	})
}
