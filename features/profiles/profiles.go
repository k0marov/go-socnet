package profiles

import (
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
)

func NewProfilesRouterImpl(db *sql.DB) func(chi.Router) {
	// db
	sqlDB, err := sql_db.NewSqlDB(db)
	if err != nil {
		log.Fatalf("Error while opening sql db as a db for profiles: %v", err)
	}
	// file storage
	avatarFileCreator := file_storage.NewAvatarFileCreator(static_file_creator.NewStaticFileCreatorImpl())
	// store
	storeDetailedProfileGetter := store.NewStoreDetailedProfileGetter(sqlDB.GetProfile)
	storeProfileUpdater := store.NewStoreProfileUpdater(sqlDB.UpdateProfile, storeDetailedProfileGetter)
	storeAvatarUpdater := store.NewStoreAvatarUpdater(avatarFileCreator, sqlDB.UpdateProfile)
	// domain
	profileUpdateValidator := validators.NewProfileUpdateValidator()
	avatarValidator := validators.NewAvatarValidator(image_decoder.ImageDecoderImpl)

	detailedProfileGetter := service.NewDetailedProfileGetter(storeDetailedProfileGetter)
	profileUpdater := service.NewProfileUpdater(profileUpdateValidator, storeProfileUpdater)
	avatarUpdater := service.NewAvatarUpdater(avatarValidator, storeAvatarUpdater)
	// handlers
	getMe := handlers.NewGetMeHandler(detailedProfileGetter)
	updateMe := handlers.NewUpdateMeHandler(profileUpdater)
	updateAvatar := handlers.NewUpdateAvatarHandler(avatarUpdater)

	return router.NewProfilesRouter(updateMe, updateAvatar, getMe)
}
