package server

import (
	"encoding/json"
	"net/http"
)

func (srv *HTTPServer) profilesMeGet(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserOrAddUnauthorized(w, r)
	if !ok {
		return
	}
	setJsonHeader(w)
	profile, err := srv.profileService.GetOrCreateDetailed(user)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}
