package server

import (
	"core/client_errors"
	core_entities "core/entities"
	"encoding/json"
	"io"
	"net/http"
	"profiles/domain/entities"
	"profiles/domain/values"

	auth "github.com/k0marov/golang-auth"
)

type AvatarData struct {
	Reader   io.Reader
	FileName string
}

type ProfileService interface {
	GetOrCreateDetailed(core_entities.User) (entities.DetailedProfile, error)
	Update(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
	UpdateAvatar(core_entities.User, AvatarData) (entities.DetailedProfile, error)
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
	} else if r.Method == http.MethodPut {
		srv.profilesMePut(w, r)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
}

func getUserOrAddUnauthorized(w http.ResponseWriter, r *http.Request) (core_entities.User, bool) {
	authUser, castSuccess := r.Context().Value(auth.UserContextKey).(auth.User) // try to cast this to User
	if !castSuccess {
		w.WriteHeader(http.StatusUnauthorized)
		return core_entities.User{}, false
	}
	return core_entities.UserFromAuth(authUser), true
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
