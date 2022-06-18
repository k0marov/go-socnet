package posts

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
)

func NewPostsRouterImpl(db *sql.DB) func(chi.Router) {
	//
	//return router.NewPostsRouter(createPost, getPost, deletePost, toggleLike)
	panic("unimplemented")
}
