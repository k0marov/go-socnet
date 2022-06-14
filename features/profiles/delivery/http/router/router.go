package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewProfilesRouter(updateMe, updateAvatar, getMe, getById, getFollowsById, toggleFollow http.HandlerFunc) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/me", getMe)
		r.Put("/me", updateMe)
		r.Put("/me/avatar", updateAvatar)

		r.Get("/{id}", getById)
		r.Get("/{id}/follows", getFollowsById)
		r.Post("/{id}/toggle-follow", toggleFollow)
	}
}
