package handlers

import (
	"github.com/k0marov/go-socnet/core/helpers/http_helpers"
	"github.com/k0marov/go-socnet/features/feed/delivery/http/responses"
	"github.com/k0marov/go-socnet/features/feed/domain/service"
	"net/http"
)

func NewFeedHandler(getFeed service.FeedGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		count := r.URL.Query().Get("count")
		posts, err := getFeed(count, caller.Id)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, responses.FeedResponse{Posts: posts})
	}
}
