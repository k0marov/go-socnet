package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/features/posts/delivery/http/handlers"
	"github.com/k0marov/socnet/features/posts/domain/service"
)

func NewPostsRouter(create service.PostCreator, getPosts service.PostsGetter, deletePost service.PostDeleter, toggleLike service.PostLikeToggler) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/", handlers.NewCreateHandler(create))
		r.Get("/?profile_id={profile_id}", handlers.NewGetListByIdHandler(getPosts))
		r.Delete("/{id}", handlers.NewDeleteHandler(deletePost))
		r.Post("/{id}/toggle-like", handlers.NewToggleLikeHandler(toggleLike))
	}
}
