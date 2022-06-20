package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewPostsRouter(create, getPosts, deletePost, toggleLike http.HandlerFunc) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/", create)
		r.Get("/", getPosts)
		r.Delete("/{id}", deletePost)
		r.Post("/{id}/toggle-like", toggleLike)
	}
}
