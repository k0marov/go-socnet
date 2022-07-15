package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewFeedRouter(feedHandler http.HandlerFunc) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", feedHandler)
	}
}
