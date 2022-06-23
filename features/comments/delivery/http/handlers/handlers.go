package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/http_helpers"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/values"
	profile_handlers "github.com/k0marov/socnet/features/profiles/delivery/http/handlers"
	"net/http"
	"time"
)

type CommentResponse struct {
	Id        string
	Author    profile_handlers.ProfileResponse
	Text      string
	CreatedAt time.Time
	Likes     int
	IsLiked   bool
	IsMine    bool
}

type CommentsResponse struct {
	Comments []CommentResponse
}

func EntityToResponse(comment entities.ContextedComment) CommentResponse {
	return CommentResponse{
		Id:        comment.Id,
		Author:    profile_handlers.EntityToResponse(comment.Author),
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		Likes:     comment.Likes,
		IsLiked:   comment.IsLiked,
		IsMine:    comment.IsMine,
	}
}

func EntitiesToResponse(comments []entities.ContextedComment) CommentsResponse {
	var commentsResp []CommentResponse
	for _, comment := range comments {
		commentsResp = append(commentsResp, EntityToResponse(comment))
	}
	return CommentsResponse{Comments: commentsResp}
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
		http_helpers.WriteJson(w, EntitiesToResponse(comments))
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
		http_helpers.WriteJson(w, EntityToResponse(createdComment))
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
