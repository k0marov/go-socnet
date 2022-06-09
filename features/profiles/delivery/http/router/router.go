package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewProfilesRouter(updateMe, updateAvatar, getMe http.HandlerFunc) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/me", getMe)
		r.Put("/me", updateMe)
		r.Put("/me/avatar", updateAvatar)
	}
}
