package posts

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/deletable"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable"
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/abstract/recommendable"
	"github.com/k0marov/go-socnet/core/general/image_decoder"
	static_store2 "github.com/k0marov/go-socnet/core/general/static_store"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/features/posts/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/posts/delivery/http/router"
	"github.com/k0marov/go-socnet/features/posts/domain/contexters"
	"github.com/k0marov/go-socnet/features/posts/domain/service"
	"github.com/k0marov/go-socnet/features/posts/domain/validators"
	"github.com/k0marov/go-socnet/features/posts/store"
	"github.com/k0marov/go-socnet/features/posts/store/file_storage"
	"github.com/k0marov/go-socnet/features/posts/store/sql_db"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
)

func NewPostRecommendable(db *sql.DB) recommendable.Recommendable {
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("error while opening sql db for posts: %v", err)
	}
	// recommendable
	recommendablePost, err := recommendable.NewRecommendable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("error while creating a Post recommendable: %v", err)
	}
	return recommendablePost
}

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

	// ownable
	ownablePost, err := ownable.NewOwnable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("error while creating a Post ownable: %v", err)
	}

	// OwnableLikeable
	ownableLikeablePost := ownable_likeable.NewOwnableLikeable(ownablePost.GetOwner, likeablePost.ToggleLike)

	// deletable
	deletablePost, err := deletable.NewDeletable(db, sqlDB.TableName, ownablePost.GetOwner)
	if err != nil {
		log.Fatalf("error while creating a Post deletable: %v", err)
	}

	// file storage
	storeImages := file_storage.NewPostImageFilesCreator(static_store2.NewStaticFileCreatorImpl())
	deleteFiles := file_storage.NewPostFilesDeleter(static_store2.NewStaticDirDeleterImpl())

	// store
	storeCreatePost := store.NewStorePostCreator(sqlDB.CreatePost, storeImages, sqlDB.AddPostImages, deletablePost.ForceDelete, deleteFiles)
	storeDeletePost := store.NewStorePostDeleter(deletablePost.ForceDelete, deleteFiles)
	storeGetPosts := store.NewStorePostsGetter(sqlDB.GetPosts, likeablePost.GetLikesCount)

	// service
	validatePost := validators.NewPostValidator(image_decoder.ImageDecoderImpl)

	// contexters
	addContext := contexters.NewPostListContextAdder(contexters.NewPostContextAdder(getContextedProfile, likeable_contexters.NewOwnLikeContextGetter(likeablePost.IsLiked)))

	createPost := service.NewPostCreator(validatePost, storeCreatePost)
	deletePost := service.NewPostDeleter(ownablePost.GetOwner, storeDeletePost)
	getPosts := service.NewPostsGetter(storeGetPosts, addContext)
	toggleLike := service.NewPostLikeToggler(ownableLikeablePost.SafeToggleLike)

	// handlers
	createPostHandler := handlers.NewCreateHandler(createPost)
	deletePostHandler := handlers.NewDeleteHandler(deletePost)
	getPostsHandler := handlers.NewGetListByIdHandler(getPosts)
	toggleLikeHandler := handlers.NewToggleLikeHandler(toggleLike)

	return router.NewPostsRouter(createPostHandler, getPostsHandler, deletePostHandler, toggleLikeHandler)
}
