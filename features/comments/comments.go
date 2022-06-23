package comments

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/core/likeable"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	"github.com/k0marov/socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/socnet/features/comments/delivery/http/router"
	"github.com/k0marov/socnet/features/comments/domain/contexters"
	"github.com/k0marov/socnet/features/comments/domain/service"
	"github.com/k0marov/socnet/features/comments/domain/validators"
	"github.com/k0marov/socnet/features/comments/store"
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
	// likeable
	likeableComment, err := likeable.NewLikeable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("error while creating comment likeable: %v", err)
	}

	// store
	storeCreateComment := store.NewCommentCreator(sqlDB.Create)
	storeGetAuthor := store.NewAuthorGetter(sqlDB.GetAuthor)
	storeGetComments := store.NewCommentsGetter(sqlDB.GetComments, likeableComment.GetLikesCount)

	// service
	validator := validators.NewCommentValidator()
	contextAdder := contexters.NewCommentListContextAdder(contexters.NewCommentContextAdder(getProfile, likeable_contexters.NewLikeableContextGetter(likeableComment.IsLiked)))

	getComments := service.NewPostCommentsGetter(storeGetComments, contextAdder)
	createComment := service.NewCommentCreator(validator, getProfile, storeCreateComment)
	toggleLike := service.NewCommentLikeToggler(storeGetAuthor, likeableComment.ToggleLike)
	// handlers
	getCommentsHandler := handlers.NewGetCommentsHandler(getComments)
	createCommentHandler := handlers.NewCreateCommentHandler(createComment)
	toggleLikeHandler := handlers.NewToggleLikeCommentHandler(toggleLike)
	return router.NewCommentsRouter(getCommentsHandler, createCommentHandler, toggleLikeHandler)
}
