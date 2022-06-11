package handlers

import (
	"core/client_errors"
	"encoding/json"
	"net/http"
	"profiles/domain/service"

	"github.com/go-chi/chi/v5"
)

func NewGetMeHandler(detailedProfileGetter service.DetailedProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		setJsonHeader(w)
		profile, err := detailedProfileGetter(user)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		json.NewEncoder(w).Encode(profile)
	})
}

func NewGetByIdHandler(profileGetter service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setJsonHeader(w)
		id := chi.URLParam(r, "id")
		if id == "" {
			throwClientError(w, client_errors.IdNotProvided)
			return
		}
		profile, err := profileGetter(id)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		json.NewEncoder(w).Encode(profile)
	})
}
