package responses

import (
	"github.com/k0marov/socnet/features/comments/domain/entities"
	profile_responses "github.com/k0marov/socnet/features/profiles/delivery/http/responses"
	"time"
)

type CommentResponse struct {
	Id        string
	Author    profile_responses.ProfileResponse
	Text      string
	CreatedAt time.Time
	Likes     int
	IsLiked   bool
	IsMine    bool
}

type CommentsResponse struct {
	Comments []CommentResponse
}

func NewCommentResponse(comment entities.ContextedComment) CommentResponse {
	return CommentResponse{
		Id:        comment.Id,
		Author:    profile_responses.NewProfileResponse(comment.Author),
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		Likes:     comment.Likes,
		IsLiked:   comment.IsLiked,
		IsMine:    comment.IsMine,
	}
}

func NewCommentListResponse(comments []entities.ContextedComment) CommentsResponse {
	var commentsResp []CommentResponse
	for _, comment := range comments {
		commentsResp = append(commentsResp, NewCommentResponse(comment))
	}
	return CommentsResponse{Comments: commentsResp}
}
