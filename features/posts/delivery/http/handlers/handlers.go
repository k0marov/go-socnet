package handlers

import (
	"net/http"
	"posts/domain/service"
)

func NewDeleteHandler(deletePost service.PostDeleter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func NewGetListByIdHandler(getPosts service.PostsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func NewToggleLikeHandler(toggleLike service.PostLikeToggler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func NewCreateNewHandler(createNew service.PostCreater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
