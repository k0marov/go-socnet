package store

import (
	"core/core_errors"
	"core/ref"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store_contracts"
	"profiles/domain/values"
)

type DBUpdateData struct {
	About      string
	AvatarPath string
}

type AvatarFileCreator = func(data ref.Ref[[]byte], belongsToUser string) (string, error)
type DBProfileUpdater = func(id string, updData DBUpdateData) error

func NewStoreAvatarUpdater(createFile AvatarFileCreator, updateDBProfile DBProfileUpdater) store_contracts.StoreAvatarUpdater {
	return func(userId string, avatar values.AvatarData) (values.AvatarPath, error) {
		path, err := createFile(avatar.Data, userId)
		if err != nil {
			return values.AvatarPath{}, fmt.Errorf("error while storing avatar file: %w", err)
		}
		avatarPath := values.AvatarPath{Path: path}
		updateData := DBUpdateData{
			AvatarPath: avatarPath.Path,
		}
		err = updateDBProfile(userId, updateData)
		if err != nil {
			return values.AvatarPath{}, fmt.Errorf("error while updating avatar path in DB: %w", err)
		}
		return avatarPath, nil
	}
}

type DBProfileCreator = func(entities.Profile) error

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store_contracts.StoreProfileCreator {
	return func(newProfile entities.Profile) error {
		return createDBProfile(newProfile)
	}
}

type DBProfileGetter = func(id string) (entities.Profile, error)

func NewStoreDetailedProfileGetter(getDBProfile DBProfileGetter) store_contracts.StoreDetailedProfileGetter {
	return func(id string) (entities.DetailedProfile, error) {
		profile, err := getDBProfile(id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.DetailedProfile{}, core_errors.ErrNotFound
			}
			return entities.DetailedProfile{}, fmt.Errorf("error while getting a profile from db: %w", err)
		}
		detailedProfile := entities.DetailedProfile{Profile: profile}
		return detailedProfile, nil
	}
}

func NewStoreProfileUpdater(updateDBProfile DBProfileUpdater, getProfile store_contracts.StoreDetailedProfileGetter) store_contracts.StoreProfileUpdater {
	return func(id string, upd values.ProfileUpdateData) (entities.DetailedProfile, error) {
		err := updateDBProfile(id, DBUpdateData{About: upd.About})
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("error while updating profile in db: %w", err)
		}
		return getProfile(id)
	}
}
