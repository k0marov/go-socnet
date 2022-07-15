package profiles

import (
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_entities"
	"github.com/k0marov/go-socnet/core/general/image_decoder"
	"github.com/k0marov/go-socnet/core/general/static_store"
	"log"

	"github.com/k0marov/go-socnet/features/profiles/domain/contexters"

	"github.com/k0marov/go-socnet/features/profiles/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/profiles/delivery/http/router"
	"github.com/k0marov/go-socnet/features/profiles/domain/service"
	"github.com/k0marov/go-socnet/features/profiles/domain/validators"
	"github.com/k0marov/go-socnet/features/profiles/store"
	"github.com/k0marov/go-socnet/features/profiles/store/file_storage"
	"github.com/k0marov/go-socnet/features/profiles/store/sql_db"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

func NewRegisterCallback(db *sqlx.DB) func(auth.User) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	// store
	storeProfileCreator := store.NewStoreProfileCreator(sqlDB.CreateProfile)
	// domain
	createProfile := service.NewProfileCreator(storeProfileCreator)
	return func(u auth.User) {
		createProfile(core_entities.UserFromAuth(u))
	}
}

func NewProfileGetterImpl(db *sqlx.DB) service.ProfileGetter {
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	likeableProfile, err := likeable.NewLikeable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("Error while creating a likeable Profile: %v", err)
	}

	addContext := contexters.NewProfileContextAdder(likeable_contexters.NewOwnLikeContextGetter(likeableProfile.IsLiked))

	getProfile := store.NewStoreProfileGetter(sqlDB.GetProfile, likeableProfile.GetLikesCount, likeableProfile.GetUserLikesCount)
	return service.NewProfileGetter(getProfile, addContext)
}

func NewProfilesRouterImpl(db *sqlx.DB) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	// likeable
	likeableProfile, err := likeable.NewLikeable(db, sqlDB.TableName)
	if err != nil {
		log.Fatalf("Error while creating a likeable Profile: %v", err)
	}

	// file storage
	avatarFileCreator := file_storage.NewAvatarFileCreator(static_store.NewStaticFileCreatorImpl())

	// store
	storeProfileGetter := store.NewStoreProfileGetter(sqlDB.GetProfile, likeableProfile.GetLikesCount, likeableProfile.GetUserLikesCount)
	storeProfileUpdater := store.NewStoreProfileUpdater(sqlDB.UpdateProfile)
	storeAvatarUpdater := store.NewStoreAvatarUpdater(avatarFileCreator, sqlDB.UpdateProfile)

	// domain
	profileUpdateValidator := validators.NewProfileUpdateValidator()
	avatarValidator := validators.NewAvatarValidator(image_decoder.ImageDecoderImpl)

	addContext := contexters.NewProfileContextAdder(likeable_contexters.NewOwnLikeContextGetter(likeableProfile.IsLiked))

	profileGetter := service.NewProfileGetter(storeProfileGetter, addContext)
	profileUpdater := service.NewProfileUpdater(profileUpdateValidator, storeProfileUpdater, profileGetter)
	avatarUpdater := service.NewAvatarUpdater(avatarValidator, storeAvatarUpdater)
	followToggler := service.NewFollowToggler(likeableProfile.ToggleLike)
	followsGetter := service.NewFollowsGetter(likeableProfile.GetUserLikes, profileGetter)

	// handlers
	getMe := handlers.NewGetMeHandler(profileGetter)
	updateMe := handlers.NewUpdateMeHandler(profileUpdater)
	updateAvatar := handlers.NewUpdateAvatarHandler(avatarUpdater)
	getFollows := handlers.NewGetFollowsHandler(followsGetter)
	getById := handlers.NewGetByIdHandler(profileGetter)
	toggleFollow := handlers.NewToggleFollowHandler(followToggler)

	return router.NewProfilesRouter(updateMe, updateAvatar, getMe, getById, getFollows, toggleFollow)
}
