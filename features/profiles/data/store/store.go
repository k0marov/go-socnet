package store

import (
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store_contracts"
	"profiles/domain/values"
)

const AvatarsDir = "avatars"

type FileCreator = func(data []byte, dir, filename string) (string, error)
type DBAvatarUpdater = func(userId string, avatarPath values.AvatarURL) error

func NewStoreAvatarUpdater(createFile FileCreator, updateDBAvatar DBAvatarUpdater) store_contracts.StoreAvatarUpdater {
	return func(userId string, avatar values.AvatarData) (values.AvatarURL, error) {
		if avatar.Data == nil {
			return values.AvatarURL{}, fmt.Errorf("provided avatar data was nil")
		}
		path, err := createFile(*avatar.Data, AvatarsDir, userId)
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

type DBProfileCreator = func(entities.DetailedProfile) error

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store_contracts.StoreProfileCreator {
	return func(newProfile entities.DetailedProfile) error {
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
