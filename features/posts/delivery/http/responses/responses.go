package responses

import (
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile_responses "github.com/k0marov/socnet/features/profiles/delivery/http/responses"
	"time"
)

type PostImageResponse struct {
	Index int
	Url   string
}

func newPostImageListResponse(images []values.PostImage) (respList []PostImageResponse) {
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
	Id        string
	Author    profile_responses.ProfileResponse
	Text      string
	CreatedAt time.Time
	Images    []PostImageResponse
	Likes     int
	IsLiked   bool
	IsMine    bool
}
type PostsResponse struct {
	Posts []PostResponse
}

func NewPostListResponse(posts []entities.ContextedPost) PostsResponse {
	var postResponses []PostResponse
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
