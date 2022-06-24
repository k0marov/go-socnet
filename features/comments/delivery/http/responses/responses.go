package responses

import (
	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	profile_responses "github.com/k0marov/go-socnet/features/profiles/delivery/http/responses"
)

type CommentResponse struct {
	Id        string                            `json:"id"`
	Author    profile_responses.ProfileResponse `json:"author"`
	Text      string                            `json:"text"`
	CreatedAt int64                             `json:"created_at"`
	Likes     int                               `json:"likes"`
	IsLiked   bool                              `json:"is_liked"`
	IsMine    bool                              `json:"is_mine"`
}

type CommentsResponse struct {
	Comments []CommentResponse `json:"comments"`
}

func NewCommentResponse(comment entities.ContextedComment) CommentResponse {
	return CommentResponse{
		Id:        comment.Id,
		Author:    profile_responses.NewProfileResponse(comment.Author),
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt.Unix(),
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
