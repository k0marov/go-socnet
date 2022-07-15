package responses

import (
	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	profile_responses "github.com/k0marov/go-socnet/features/profiles/delivery/http/responses"
)

type PostImageResponse struct {
	Index int    `json:"index"`
	Url   string `json:"url"`
}

func newPostImageListResponse(images []values.PostImage) []PostImageResponse {
	respList := make([]PostImageResponse, 0)
	for _, img := range images {
		resp := PostImageResponse{
			Index: img.Index,
			Url:   img.URL,
		}
		respList = append(respList, resp)
	}
	return respList
}

type PostResponse struct {
	Id        string                            `json:"id"`
	Author    profile_responses.ProfileResponse `json:"author"`
	Text      string                            `json:"text"`
	CreatedAt int64                             `json:"created_at"`
	Images    []PostImageResponse               `json:"images"`
	Likes     int                               `json:"likes"`
	IsLiked   bool                              `json:"is_liked"`
	IsMine    bool                              `json:"is_mine"`
}
type PostsResponse struct {
	Posts []PostResponse `json:"posts"`
}

func NewPostListResponse(posts []entities.ContextedPost) PostsResponse {
	postResponses := make([]PostResponse, 0)
	for _, post := range posts {
		resp := PostResponse{
			Id:        post.Id,
			Author:    profile_responses.NewProfileResponse(post.Author),
			Text:      post.Text,
			CreatedAt: post.CreatedAt,
			Images:    newPostImageListResponse(post.Images),
			Likes:     post.Likes,
			IsLiked:   post.IsLiked,
			IsMine:    post.IsMine,
		}
		postResponses = append(postResponses, resp)
	}
	return PostsResponse{
		Posts: postResponses,
	}
}
