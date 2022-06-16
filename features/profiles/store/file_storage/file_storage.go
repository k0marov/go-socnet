package file_storage

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/profiles/store"

	"github.com/k0marov/socnet/core/ref"
	"github.com/k0marov/socnet/core/static_store"
)

const ProfilePrefix = "profile_"
const AvatarFileName = "avatar"

func NewAvatarFileCreator(createFile static_store.StaticFileCreator) store.AvatarFileCreator {
	return func(data ref.Ref[[]byte], belongsToUser core_values.UserId) (core_values.StaticFilePath, error) {
		return createFile(data, ProfilePrefix+belongsToUser, AvatarFileName)
	}
}
