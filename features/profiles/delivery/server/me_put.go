package server

import (
	"core/client_errors"
	"encoding/json"
	"io"
	"net/http"
	"profiles/domain/values"
	"strings"
)

func (srv *HTTPServer) profilesMePut(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSuffix(r.URL.String(), "/")
	if strings.HasSuffix(url, "avatar") {
		srv.profilesMeUpdateAvatar(w, r)
	} else {
		srv.profilesMeUpdate(w, r)
	}
}

func (srv *HTTPServer) profilesMeUpdate(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserOrAddUnauthorized(w, r)
	if !ok {
		return
	}

	setJsonHeader(w)

	var updateData values.ProfileUpdateData
	json.NewDecoder(r.Body).Decode(&updateData)

	updatedProfile, err := srv.profileService.Update(user, updateData)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(updatedProfile)
}

const MaxFileSize = 3 << 20 // 3 MB

func (srv *HTTPServer) profilesMeUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserOrAddUnauthorized(w, r)
	if !ok {
		return
	}

	avatarData, clientError := _parseAvatar(r)
	if clientError != client_errors.NoError {
		throwClientError(w, clientError)
		return
	}

	setJsonHeader(w)
	profile, err := srv.profileService.UpdateAvatar(user, avatarData)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func _parseAvatar(r *http.Request) (values.AvatarData, client_errors.ClientError) {
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		return values.AvatarData{}, client_errors.BodyIsNotMultipartForm
	}
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		return values.AvatarData{}, client_errors.AvatarNotProvidedError
	}
	avatarData, _ := io.ReadAll(file) // ignore error here, because later empty avatarData will throw NonImageAvatar
	file.Close()
	return values.AvatarData{Data: &avatarData, FileName: fileHeader.Filename}, client_errors.NoError
}
