package file_storage

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/core_values/ref"
	"github.com/k0marov/go-socnet/core/general/static_store"
	"github.com/k0marov/go-socnet/features/profiles/store"
)

const ProfilePrefix = "profile_"
const AvatarFileName = "avatar"

func NewAvatarFileCreator(createFile static_store.StaticFileCreator) store.AvatarFileCreator {
	return func(data ref.Ref[[]byte], belongsToUser core_values.UserId) (core_values.StaticPath, error) {
		return createFile(data, ProfilePrefix+belongsToUser, AvatarFileName)
	}
}

func GetProfileDir(user core_values.UserId) core_values.StaticPath {
	return ProfilePrefix + user
}
