package handlers

import (
	"core/client_errors"
	"core/ref"
	"encoding/json"
	"io"
	"net/http"
	"profiles/domain/service"
	"profiles/domain/values"
)

func NewUpdateMeHandler(profileUpdater service.ProfileUpdater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}

		setJsonHeader(w)

		var updateData values.ProfileUpdateData
		err := json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			throwClientError(w, client_errors.InvalidJsonError)
			return
		}

		updatedProfile, err := profileUpdater(user, updateData)
		if err != nil {
			handleServiceError(w, err)
			return
		}

		json.NewEncoder(w).Encode(updatedProfile)
	})
}

func NewUpdateAvatarHandler(avatarUpdater service.AvatarUpdater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}

		avatarData, ok := _parseAvatar(r)
		if !ok {
			throwClientError(w, client_errors.AvatarNotProvidedError)
			return
		}

		setJsonHeader(w)
		url, err := avatarUpdater(user, avatarData)
		if err != nil {
			handleServiceError(w, err)
			return
		}

		json.NewEncoder(w).Encode(url)
	})
}

const MaxFileSize = 3 << 20 // 3 MB
func _parseAvatar(r *http.Request) (values.AvatarData, bool) {
	file, _, err := r.FormFile("avatar")
	if err != nil {
		return values.AvatarData{}, false
	}
	defer file.Close()
	avatarData, err := io.ReadAll(file)
	if err != nil {
		return values.AvatarData{}, false
	}
	dataRef, err := ref.NewRef(&avatarData)
	if err != nil {
		return values.AvatarData{}, false
	}
	return values.AvatarData{Data: dataRef}, true
}
