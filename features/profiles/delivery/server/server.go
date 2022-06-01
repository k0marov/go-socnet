package server

import (
	"core/client_errors"
	"encoding/json"
	"net/http"
	"profiles/domain/entities"

	auth "github.com/k0marov/golang-auth"
)

type ProfileService interface {
	GetOrCreate(auth.User) (entities.Profile, error)
}

type HTTPServer struct {
	profileService ProfileService
}

func NewHTTPServer(service ProfileService) *HTTPServer {
	return &HTTPServer{profileService: service}
}

func (srv *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, castSuccess := r.Context().Value("User").(auth.User) // try to cast this to User
	if !castSuccess {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("contentType", "application/json")
	profile, err := srv.profileService.GetOrCreate(user)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func handleServiceError(w http.ResponseWriter, err error) {
	clientError, isClientError := err.(client_errors.ClientError)
	if isClientError {
		errorJson, _ := json.Marshal(clientError)
		http.Error(w, string(errorJson), http.StatusBadRequest)
	} else {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
