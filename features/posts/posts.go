package posts

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/core/image_decoder"
	"github.com/k0marov/socnet/core/likeable"
	"github.com/k0marov/socnet/core/static_store"
	"github.com/k0marov/socnet/features/posts/delivery/http/handlers"
	"github.com/k0marov/socnet/features/posts/delivery/http/router"
	"github.com/k0marov/socnet/features/posts/domain/service"
	"github.com/k0marov/socnet/features/posts/domain/validators"
	"github.com/k0marov/socnet/features/posts/store"
	"github.com/k0marov/socnet/features/posts/store/file_storage"
	"github.com/k0marov/socnet/features/posts/store/sql_db"
	profile_service "github.com/k0marov/socnet/features/profiles/domain/service"
	"log"
)

func NewPostsRouterImpl(db *sql.DB, getContextedProfile profile_service.ProfileGetter) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("error while opening sql db for posts: %v", err)
	}

	// likeable
	likeablePost, err := likeable.NewLikeable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("error while creating a Post likeable: %v", err)
	}

	// file storage
	storeImages := file_storage.NewPostImageFilesCreator(static_store.NewStaticFileCreatorImpl())
	deleteFiles := file_storage.NewPostFilesDeleter(static_store.NewStaticDirDeleterImpl())

	// store
	storeCreatePost := store.NewStorePostCreator(sqlDB.CreatePost, storeImages, sqlDB.AddPostImages, sqlDB.DeletePost, deleteFiles)
	storeGetAuthor := store.NewStoreAuthorGetter(sqlDB.GetAuthor)
	storeDeletePost := store.NewStorePostDeleter(sqlDB.DeletePost, deleteFiles)
	storeGetPosts := store.NewStorePostsGetter(sqlDB.GetPosts)

	// service
	validatePost := validators.NewPostValidator(image_decoder.ImageDecoderImpl)

	createPost := service.NewPostCreator(validatePost, storeCreatePost)
	deletePost := service.NewPostDeleter(storeGetAuthor, storeDeletePost)
	getPosts := service.NewPostsGetter(getContextedProfile, storeGetPosts, likeablePost.GetLikesCount, likeablePost.IsLiked)
	toggleLike := service.NewPostLikeToggler(storeGetAuthor, likeablePost.ToggleLike)

	// handlers
	createPostHandler := handlers.NewCreateHandler(createPost)
	deletePostHandler := handlers.NewDeleteHandler(deletePost)
	getPostsHandler := handlers.NewGetListByIdHandler(getPosts)
	toggleLikeHandler := handlers.NewToggleLikeHandler(toggleLike)

	return router.NewPostsRouter(createPostHandler, getPostsHandler, deletePostHandler, toggleLikeHandler)
}
