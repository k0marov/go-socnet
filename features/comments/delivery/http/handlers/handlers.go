package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/http_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"net/http"
)

//type CommentResponse struct {
//	Id string
//	Author ProfileResponse
//
//}
type CommentsResponse struct {
	Comments []entities.ContextedComment
}
type NewCommentRequest struct {
	Text string
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
		http_helpers.WriteJson(w, CommentsResponse{comments})
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
		http_helpers.WriteJson(w, createdComment)
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
