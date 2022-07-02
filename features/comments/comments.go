package comments

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/deletable"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable"
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/router"
	"github.com/k0marov/go-socnet/features/comments/domain/contexters"
	"github.com/k0marov/go-socnet/features/comments/domain/service"
	"github.com/k0marov/go-socnet/features/comments/domain/validators"
	"github.com/k0marov/go-socnet/features/comments/store"
	"github.com/k0marov/go-socnet/features/comments/store/sql_db"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
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
	// ownable
	ownableComment, err := ownable.NewOwnable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("error while creating comment ownable: %v", err)
	}
	// ownable-likeable
	ownableLikeableComment := ownable_likeable.NewOwnableLikeable(ownableComment.GetOwner, likeableComment.ToggleLike)

	// deletable
	deletableComment, err := deletable.NewDeletable(db, sqlDB.TableName, ownableComment.GetOwner)
	if err != nil {
		log.Fatalf("error while creating comment deletable: %v", err)
	}

	// store
	storeCreateComment := store.NewCommentCreator(sqlDB.Create)
	storeGetComments := store.NewCommentsGetter(sqlDB.GetComments, likeableComment.GetLikesCount)

	// service
	validator := validators.NewCommentValidator()
	contextAdder := contexters.NewCommentListContextAdder(contexters.NewCommentContextAdder(getProfile, likeable_contexters.NewOwnLikeContextGetter(likeableComment.IsLiked)))

	getComments := service.NewPostCommentsGetter(storeGetComments, contextAdder)
	createComment := service.NewCommentCreator(validator, getProfile, storeCreateComment)
	toggleLike := service.NewCommentLikeToggler(ownableLikeableComment.SafeToggleLike)
	delete := service.NewCommentDeleter(deletableComment.Delete)
	// handlers
	getCommentsHandler := handlers.NewGetCommentsHandler(getComments)
	createCommentHandler := handlers.NewCreateCommentHandler(createComment)
	toggleLikeHandler := handlers.NewToggleLikeCommentHandler(toggleLike)
	deleteHandler := handlers.NewDeleteCommentHandler(delete)
	return router.NewCommentsRouter(getCommentsHandler, createCommentHandler, toggleLikeHandler, deleteHandler)
}
