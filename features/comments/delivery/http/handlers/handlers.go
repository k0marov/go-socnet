package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/http_helpers"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/responses"
	"github.com/k0marov/go-socnet/features/comments/domain/service"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
)

type NewCommentRequest struct {
	Text string `json:"text"`
}

func NewGetCommentsHandler(getComments service.PostCommentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := r.URL.Query().Get("post_id")
		if postId == "" {
			http_helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		comments, err := getComments(postId, caller.Id)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, responses.NewCommentListResponse(comments))
	}
}
func NewCreateCommentHandler(createComment service.CommentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := r.URL.Query().Get("post_id")
		if postId == "" {
			http_helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		var commentData NewCommentRequest
		json.NewDecoder(r.Body).Decode(&commentData)
		newComment := values.NewCommentValue{
			Text:   commentData.Text,
			Author: caller.Id,
			Post:   postId,
		}
		createdComment, err := createComment(newComment)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, responses.NewCommentResponse(createdComment))
	}
}

func NewToggleLikeCommentHandler(toggleLike service.CommentLikeToggler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		commentId := chi.URLParam(r, "id")
		if commentId == "" {
			http_helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		err := toggleLike(commentId, caller.Id)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
	}
}
