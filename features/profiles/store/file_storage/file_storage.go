package file_storage

import (
	"core/ref"
	"core/static_file_creator"
	"profiles/data/store"
)

const UserPrefix = "user_"
const AvatarFileName = "avatar"

func NewAvatarFileCreator(createFile static_file_creator.StaticFileCreator) store.AvatarFileCreator {
	return func(data ref.Ref[[]byte], belongsToUser string) (string, error) {
		return createFile(data, UserPrefix+belongsToUser, AvatarFileName)
	}
}
