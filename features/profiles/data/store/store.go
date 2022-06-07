package store

import (
	"fmt"
	"profiles/domain/store_contracts"
	"profiles/domain/values"
)

const AvatarsDir = "avatars"

type FileCreator = func(data *[]byte, dir, filename string) (string, error)
type DBAvatarUpdater = func(userId string, avatarPath values.AvatarURL) error

func NewStoreAvatarUpdater(createFile FileCreator, updateDBAvatar DBAvatarUpdater) store_contracts.StoreAvatarUpdater {
	return func(userId string, avatar values.AvatarData) (values.AvatarURL, error) {
		path, err := createFile(avatar.Data, AvatarsDir, userId)
		if err != nil {
			return values.AvatarURL{}, fmt.Errorf("error while storing avatar file: %w", err)
		}
		avatarUrl := values.AvatarURL{Url: path}
		err = updateDBAvatar(userId, avatarUrl)
		if err != nil {
			return values.AvatarURL{}, fmt.Errorf("error while updating avatar path in DB: %w", err)
		}
		return avatarUrl, nil
	}
}
