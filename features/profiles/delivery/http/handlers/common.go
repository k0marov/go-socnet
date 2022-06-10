package handlers

import (
	"core/client_errors"
	core_entities "core/entities"
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/k0marov/golang-auth"
)

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
		throwClientError(w, clientError)
	} else {
		log.Printf("Error while serving request: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func throwClientError(w http.ResponseWriter, clientError client_errors.ClientError) {
	setJsonHeader(w)
	errorJson, _ := json.Marshal(clientError)
	http.Error(w, string(errorJson), clientError.HTTPCode)
}
