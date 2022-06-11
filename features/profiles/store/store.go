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

type AvatarFileCreator = func(data ref.Ref[[]byte], belongsToUser values.UserId) (string, error)
type DBProfileUpdater = func(id values.UserId, updData DBUpdateData) error

func NewStoreAvatarUpdater(createFile AvatarFileCreator, updateDBProfile DBProfileUpdater) store_contracts.StoreAvatarUpdater {
	return func(userId values.UserId, avatar values.AvatarData) (values.AvatarPath, error) {
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

type DBProfileGetter = func(id values.UserId) (entities.Profile, error)

func NewStoreDetailedProfileGetter(getDBProfile DBProfileGetter) store_contracts.StoreDetailedProfileGetter {
	return func(id values.UserId) (entities.DetailedProfile, error) {
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
	return func(id values.UserId, upd values.ProfileUpdateData) (entities.DetailedProfile, error) {
		err := updateDBProfile(id, DBUpdateData{About: upd.About})
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("error while updating profile in db: %w", err)
		}
		return getProfile(id)
	}
}

type DBProfileCreator = func(entities.Profile) error

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store_contracts.StoreProfileCreator {
	return createDBProfile
}

func NewStoreProfileGetter(getDBProfile DBProfileGetter) store_contracts.StoreProfileGetter {
	return getDBProfile
}

type DBFollowsGetter = func(id values.UserId) ([]entities.Profile, error)

func NewStoreFollowsGetter(dbFollowsGetter DBFollowsGetter) store_contracts.StoreFollowsGetter {
	return dbFollowsGetter
}

type DBFollowChecker = func(target, follower values.UserId) (bool, error)

func NewStoreFollowChecker(dbFollowChecker DBFollowChecker) store_contracts.StoreFollowChecker {
	return dbFollowChecker
}

type DBFollower = func(target, follower values.UserId) error

func NewStoreFollower(dbFollower DBFollower) store_contracts.StoreFollower {
	return dbFollower
}

type DBUnfollower = func(target, unfollower values.UserId) error

func NewStoreUnfollower(dbUnfollower DBUnfollower) store_contracts.StoreUnfollower {
	return dbUnfollower
}
