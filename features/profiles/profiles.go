package profiles

import (
	"database/sql"
	"log"

	"github.com/k0marov/socnet/features/profiles/delivery/http/handlers"
	"github.com/k0marov/socnet/features/profiles/delivery/http/router"
	"github.com/k0marov/socnet/features/profiles/domain/service"
	"github.com/k0marov/socnet/features/profiles/domain/validators"
	"github.com/k0marov/socnet/features/profiles/store"
	"github.com/k0marov/socnet/features/profiles/store/file_storage"
	"github.com/k0marov/socnet/features/profiles/store/sql_db"

	"github.com/k0marov/socnet/core/core_entities"
	"github.com/k0marov/socnet/core/image_decoder"
	"github.com/k0marov/socnet/core/static_store"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

func NewRegisterCallback(db *sql.DB) func(auth.User) {
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

func NewProfilesRouterImpl(db *sql.DB) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	// file storage
	avatarFileCreator := file_storage.NewAvatarFileCreator(static_store.NewStaticFileCreatorImpl())
	// store
	storeProfileUpdater := store.NewStoreProfileUpdater(sqlDB.UpdateProfile, sqlDB.GetProfile)
	storeAvatarUpdater := store.NewStoreAvatarUpdater(avatarFileCreator, sqlDB.UpdateProfile)
	storeFollowsGetter := store.NewStoreFollowsGetter(sqlDB.GetFollows)
	storeProfileGetter := store.NewStoreProfileGetter(sqlDB.GetProfile)
	storeFollowChecker := store.NewStoreFollowChecker(sqlDB.IsFollowing)
	storeFollower := store.NewStoreFollower(sqlDB.Follow)
	storeUnfollower := store.NewStoreUnfollower(sqlDB.Unfollow)
	// domain
	profileUpdateValidator := validators.NewProfileUpdateValidator()
	avatarValidator := validators.NewAvatarValidator(image_decoder.ImageDecoderImpl)

	profileUpdater := service.NewProfileUpdater(profileUpdateValidator, storeProfileUpdater)
	avatarUpdater := service.NewAvatarUpdater(avatarValidator, storeAvatarUpdater)
	followsGetter := service.NewFollowsGetter(storeFollowsGetter)
	profileGetter := service.NewProfileGetter(storeProfileGetter, storeFollowChecker)
	followToggler := service.NewFollowToggler(storeFollowChecker, storeFollower, storeUnfollower)
	// handlers
	getMe := handlers.NewGetMeHandler(profileGetter)
	updateMe := handlers.NewUpdateMeHandler(profileUpdater)
	updateAvatar := handlers.NewUpdateAvatarHandler(avatarUpdater)
	getFollows := handlers.NewGetFollowsHandler(followsGetter)
	getById := handlers.NewGetByIdHandler(profileGetter)
	toggleFollow := handlers.NewToggleFollowHandler(followToggler)

	return router.NewProfilesRouter(updateMe, updateAvatar, getMe, getById, getFollows, toggleFollow)
}
