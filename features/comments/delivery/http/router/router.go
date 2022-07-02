package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewCommentsRouter(getComments, createComment, toggleLike, delete http.HandlerFunc) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", getComments)
		r.Post("/", createComment)
		r.Post("/{id}/toggle-like", toggleLike)
		r.Delete("/{id}", delete)
	}
}
