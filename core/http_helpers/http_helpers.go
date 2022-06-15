package http_helpers

import (
	"core/client_errors"
	"core/core_values"
	core_entities "core/entities"
	"core/ref"
	"encoding/json"
	"io"
	"log"
	"net/http"

	auth "github.com/k0marov/golang-auth"
)

func SetJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
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
	SetJsonHeader(w)
	errorJson, _ := json.Marshal(clientError)
	http.Error(w, string(errorJson), clientError.HTTPCode)
}

const MaxFileSize = 3 << 20 // 3 MB
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
