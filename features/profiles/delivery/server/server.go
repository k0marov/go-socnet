package server

import (
	"core/client_errors"
	"encoding/json"
	"io"
	"net/http"
	"profiles/domain/entities"
	"profiles/domain/values"
	"strings"

	auth "github.com/k0marov/golang-auth"
)

type AvatarData struct {
	Reader   io.Reader
	FileName string
}

type ProfileService interface {
	GetOrCreate(auth.User) (entities.Profile, error)
	Update(auth.User, values.ProfileUpdateData) (entities.Profile, error)
	UpdateAvatar(auth.User, AvatarData) (entities.Profile, error)
}

type HTTPServer struct {
	profileService ProfileService
}

func NewHTTPServer(service ProfileService) *HTTPServer {
	return &HTTPServer{profileService: service}
}

func (srv *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		srv.profilesMeGet(w, r)
	} else if r.Method == http.MethodPost {
		srv.profilesMePost(w, r)
	}
}

func (srv *HTTPServer) profilesMePost(w http.ResponseWriter, r *http.Request) {
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

const MaxFileSize = 2.5 * 1000 * 1000

func (srv *HTTPServer) profilesMeUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserOrAddUnauthorized(w, r)
	if !ok {
		return
	}

	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		throwClientError(w, client_errors.BodyIsNotMultipartForm, http.StatusBadRequest)
	}
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		throwClientError(w, client_errors.AvatarNotProvidedError, http.StatusBadRequest)
		return
	}
	if fileHeader.Size > MaxFileSize {
		throwClientError(w, client_errors.AvatarTooBigError, http.StatusBadRequest)
		return
	}

	setJsonHeader(w)
	profile, err := srv.profileService.UpdateAvatar(user, AvatarData{file, fileHeader.Filename})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (srv *HTTPServer) profilesMeGet(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserOrAddUnauthorized(w, r)
	if !ok {
		return
	}
	setJsonHeader(w)
	profile, err := srv.profileService.GetOrCreate(user)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
}

func getUserOrAddUnauthorized(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	user, castSuccess := r.Context().Value(auth.UserContextKey).(auth.User) // try to cast this to User
	if !castSuccess {
		w.WriteHeader(http.StatusUnauthorized)
		return auth.User{}, false
	}
	return user, true
}

func handleServiceError(w http.ResponseWriter, err error) {
	clientError, isClientError := err.(client_errors.ClientError)
	if isClientError {
		throwClientError(w, clientError, http.StatusBadRequest)
	} else {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func throwClientError(w http.ResponseWriter, clientError client_errors.ClientError, statusCode int) {
	setJsonHeader(w)
	errorJson, _ := json.Marshal(clientError)
	http.Error(w, string(errorJson), statusCode)
}
