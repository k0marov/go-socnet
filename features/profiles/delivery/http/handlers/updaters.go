package handlers

import (
	"core/client_errors"
	helpers "core/http_helpers"
	"encoding/json"
	"net/http"
	"profiles/domain/service"
	"profiles/domain/values"

	"github.com/go-chi/chi/v5"
)

func NewUpdateMeHandler(profileUpdater service.ProfileUpdater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}

		helpers.SetJsonHeader(w)

		var updateData values.ProfileUpdateData
		err := json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			helpers.ThrowClientError(w, client_errors.InvalidJsonError)
			return
		}

		updatedProfile, err := profileUpdater(user, updateData)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}

		json.NewEncoder(w).Encode(updatedProfile)
	})
}

func NewToggleFollowHandler(followToggler service.FollowToggler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		follower, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}

		targetId := chi.URLParam(r, "id")
		if targetId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}

		err := followToggler(targetId, follower.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
	})
}

func NewUpdateAvatarHandler(avatarUpdater service.AvatarUpdater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}

		avatarFileData, ok := helpers.ParseFile(r, "avatar")
		if !ok {
			helpers.ThrowClientError(w, client_errors.AvatarNotProvidedError)
			return
		}
		avatarData := values.AvatarData{Data: avatarFileData}

		helpers.SetJsonHeader(w)
		url, err := avatarUpdater(user, avatarData)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}

		json.NewEncoder(w).Encode(url)
	})
}
