package handlers

import (
	"encoding/json"
	"net/http"
	"profiles/domain/service"
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
