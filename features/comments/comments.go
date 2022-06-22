package comments

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/socnet/features/comments/delivery/http/router"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/validators"
	"github.com/k0marov/socnet/features/comments/store/sql_db"
	profile_service "github.com/k0marov/socnet/features/profiles/domain/service"
	"log"
)

func NewCommentsRouterImpl(db *sql.DB, getProfile profile_service.ProfileGetter) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("error while opening sql db for comments: %v", err)
	}
	// service
	validator := validators.NewCommentValidator()
	getComments := service.NewPostCommentsGetter(sqlDB.GetComments, getProfile, sqlDB.IsLiked)
	createComment := service.NewCommentCreator(validator, getProfile, sqlDB.Create)
	toggleLike := service.NewCommentLikeToggler(sqlDB.GetAuthor, sqlDB.IsLiked, sqlDB.Like, sqlDB.Unlike)
	// handlers
	getCommentsHandler := handlers.NewGetCommentsHandler(getComments)
	createCommentHandler := handlers.NewCreateCommentHandler(createComment)
	toggleLikeHandler := handlers.NewToggleLikeCommentHandler(toggleLike)
	return router.NewCommentsRouter(getCommentsHandler, createCommentHandler, toggleLikeHandler)
}
