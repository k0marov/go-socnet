package profiles

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

// func NewRegisterCallback(db *sql.DB) func(auth.User) {
// 	// db
// 	sqlDB, err := sql_db.NewSqlDB(db)
// 	if err != nil {
// 		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
// 	}
// 	// store
// 	storeProfileCreator := store.NewStoreProfileCreator(sqlDB.CreateProfile)
// 	// domain
// 	createProfile := service.NewProfileCreator(storeProfileCreator)
// 	return func(u auth.User) {
// 		createProfile(entities.UserFromAuth(u))
// 	}
// }

func NewProfilesRouterImpl(db *sql.DB) func(chi.Router) {
	panic("unimplemented")

	// // db
	// sqlDB, err := sql_db.NewSqlDB(db)
	// if err != nil {
	// 	log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	// }
	// // file storage
	// avatarFileCreator := file_storage.NewAvatarFileCreator(static_file_creator.NewStaticFileCreatorImpl())
	// // store
	// storeDetailedProfileGetter := store.NewStoreDetailedProfileGetter(sqlDB.GetProfile)
	// storeProfileUpdater := store.NewStoreProfileUpdater(sqlDB.UpdateProfile, storeDetailedProfileGetter)
	// storeAvatarUpdater := store.NewStoreAvatarUpdater(avatarFileCreator, sqlDB.UpdateProfile)
	// // domain
	// profileUpdateValidator := validators.NewProfileUpdateValidator()
	// avatarValidator := validators.NewAvatarValidator(image_decoder.ImageDecoderImpl)

	// detailedProfileGetter := service.NewDetailedProfileGetter(storeDetailedProfileGetter)
	// profileUpdater := service.NewProfileUpdater(profileUpdateValidator, storeProfileUpdater)
	// avatarUpdater := service.NewAvatarUpdater(avatarValidator, storeAvatarUpdater)
	// // handlers
	// getMe := handlers.NewGetMeHandler(detailedProfileGetter)
	// updateMe := handlers.NewUpdateMeHandler(profileUpdater)
	// updateAvatar := handlers.NewUpdateAvatarHandler(avatarUpdater)

	// return router.NewProfilesRouter(updateMe, updateAvatar, getMe)
}
