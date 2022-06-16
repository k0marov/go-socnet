package profiles

import (
	"core/entities"
	"core/image_decoder"
	"core/static_file_creator"
	"database/sql"
	"log"
	"profiles/delivery/http/handlers"
	"profiles/delivery/http/router"
	"profiles/domain/service"
	"profiles/domain/validators"
	"profiles/store"
	"profiles/store/file_storage"
	"profiles/store/sql_db"

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
		createProfile(entities.UserFromAuth(u))
	}
}

func NewProfilesRouterImpl(db *sql.DB) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	// file storage
	avatarFileCreator := file_storage.NewAvatarFileCreator(static_file_creator.NewStaticFileCreatorImpl())
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
