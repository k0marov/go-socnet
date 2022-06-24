package http_helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/k0marov/go-socnet/core/client_errors"
	core_entities "github.com/k0marov/go-socnet/core/core_entities"
	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/ref"

	auth "github.com/k0marov/golang-auth"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
}

func WriteJson(w http.ResponseWriter, obj any) {
	setJsonHeader(w)
	json.NewEncoder(w).Encode(obj)
}

func GetUserOrAddUnauthorized(w http.ResponseWriter, r *http.Request) (core_entities.User, bool) {
	authUser, castSuccess := r.Context().Value(auth.UserContextKey).(auth.User) // try to cast this to User
	if !castSuccess {
		w.WriteHeader(http.StatusUnauthorized)
		return core_entities.User{}, false
	}
	return core_entities.UserFromAuth(authUser), true
}

func HandleServiceError(w http.ResponseWriter, err error) {
	clientError, isClientError := err.(client_errors.ClientError)
	if isClientError {
		ThrowClientError(w, clientError)
	} else {
		log.Printf("Error while serving request: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func ThrowClientError(w http.ResponseWriter, clientError client_errors.ClientError) {
	setJsonHeader(w)
	errorJson, _ := json.Marshal(clientError)
	http.Error(w, string(errorJson), clientError.HTTPCode)
}

func ParseFile(r *http.Request, field string) (core_values.FileData, bool) {
	file, _, err := r.FormFile(field)
	if err != nil {
		return core_values.FileData{}, false
	}
	defer file.Close()
	avatarData, err := io.ReadAll(file)
	if err != nil {
		return core_values.FileData{}, false
	}
	dataRef, err := ref.NewRef(&avatarData)
	if err != nil {
		return core_values.FileData{}, false
	}
	return dataRef, true
}
