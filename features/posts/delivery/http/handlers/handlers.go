package handlers

import (
	profile_handlers "github.com/k0marov/socnet/features/profiles/delivery/http/handlers"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/service"
	"github.com/k0marov/socnet/features/posts/domain/values"

	"github.com/k0marov/socnet/core/client_errors"
	helpers "github.com/k0marov/socnet/core/http_helpers"

	"github.com/go-chi/chi/v5"
)

func NewDeleteHandler(deletePost service.PostDeleter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := chi.URLParam(r, "id")
		if postId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		err := deletePost(postId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
	})
}

type PostImageResponse struct {
	Index int
	Url   string
}

func PostImagesToResponse(images []values.PostImage) (respList []PostImageResponse) {
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
	Author    profile_handlers.ProfileResponse
	Text      string
	CreatedAt time.Time
	Images    []PostImageResponse
	Likes     int
	IsLiked   bool
	IsMine    bool
}

func PostsToResponse(posts []entities.ContextedPost) PostsResponse {
	var postResponses []PostResponse
	for _, post := range posts {
		resp := PostResponse{
			Id:        post.Id,
			Author:    profile_handlers.EntityToResponse(post.Author),
			Text:      post.Text,
			CreatedAt: post.CreatedAt,
			Images:    PostImagesToResponse(post.Images),
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

type PostsResponse struct {
	Posts []PostResponse
}

func NewGetListByIdHandler(getPosts service.PostsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		profileId := r.URL.Query().Get("profile_id")
		if profileId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		posts, err := getPosts(profileId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, PostsToResponse(posts))
	})
}

func NewToggleLikeHandler(toggleLike service.PostLikeToggler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		postId := chi.URLParam(r, "id")
		if postId == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		err := toggleLike(postId, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
		}
	})
}

func NewCreateHandler(createPost service.PostCreator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		newPost := values.NewPostData{
			Author: user.Id,
			Text:   r.FormValue("text"),
			Images: parseImages(r),
		}
		log.Printf("%+v", newPost)
		err := createPost(newPost)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
	})
}

func parseImages(r *http.Request) []values.PostImageFile {
	images := []values.PostImageFile{}
	for i := 1; ; i++ {
		file, ok := helpers.ParseFile(r, "image_"+strconv.Itoa(i))
		if !ok {
			return images
		}
		image := values.PostImageFile{File: file, Index: i}
		images = append(images, image)
	}
}
