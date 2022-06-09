package store

import (
	"core/ref"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store_contracts"
	"profiles/domain/values"
)

type AvatarFileCreator = func(data ref.Ref[[]byte], belongsToUser string) (string, error)
type DBAvatarUpdater = func(userId string, avatarPath values.AvatarURL) error

func NewStoreAvatarUpdater(createFile AvatarFileCreator, updateDBAvatar DBAvatarUpdater) store_contracts.StoreAvatarUpdater {
	return func(userId string, avatar values.AvatarData) (values.AvatarURL, error) {
		path, err := createFile(avatar.Data, userId)
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

type DBProfileCreator = func(values.NewProfile) error

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store_contracts.StoreProfileCreator {
	return func(newProfile values.NewProfile) error {
		return createDBProfile(newProfile)
	}
}

type DBProfileGetter = func(id string) (entities.Profile, error)

func NewStoreDetailedProfileGetter(getDBProfile DBProfileGetter) store_contracts.StoreDetailedProfileGetter {
	return func(id string) (entities.DetailedProfile, error) {
		profile, err := getDBProfile(id)
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("error while getting a profile from db: %w", err)
		}
		detailedProfile := entities.DetailedProfile{Profile: profile}
		return detailedProfile, nil
	}
}

type DBProfileUpdater = func(id string, updData values.ProfileUpdateData) error

func NewStoreProfileUpdater(updateDBProfile DBProfileUpdater, getProfile store_contracts.StoreDetailedProfileGetter) store_contracts.StoreProfileUpdater {
	return func(id string, upd values.ProfileUpdateData) (entities.DetailedProfile, error) {
		err := updateDBProfile(id, upd)
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("error while updating profile in db: %w", err)
		}
		return getProfile(id)
	}
}
