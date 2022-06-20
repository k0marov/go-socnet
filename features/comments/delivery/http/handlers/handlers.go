package handlers

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/http_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"net/http"
)

type CommentsResponse struct {
	Comments []entities.Comment
}

func NewGetCommentsHandler(getComments service.PostCommentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId := r.URL.Query().Get("post_id")
		if postId == "" {
			http_helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		comments, err := getComments(postId)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, CommentsResponse{comments})
	}
}
func NewAddCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func NewToggleLikeCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
