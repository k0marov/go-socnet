package file_storage

import (
	"github.com/k0marov/socnet/features/profiles/store"

	"github.com/k0marov/socnet/core/ref"
	"github.com/k0marov/socnet/core/static_store"
)

const UserPrefix = "user_"
const AvatarFileName = "avatar"

func NewAvatarFileCreator(createFile static_store.StaticFileCreator) store.AvatarFileCreator {
	return func(data ref.Ref[[]byte], belongsToUser string) (string, error) {
		return createFile(data, UserPrefix+belongsToUser, AvatarFileName)
	}
}
