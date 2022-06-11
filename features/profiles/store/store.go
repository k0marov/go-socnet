package store

import (
	"core/core_errors"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store_contracts"
	"profiles/domain/values"
)

type DBUpdateData struct {
	About      string
	AvatarPath string
}

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

func NewStoreDetailedProfileGetter(getDBProfile DBProfileGetter, getDBFollows DBFollowsGetter) store_contracts.StoreDetailedProfileGetter {
	return func(id values.UserId) (entities.DetailedProfile, error) {
		profile, err := getDBProfile(id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.DetailedProfile{}, core_errors.ErrNotFound
			}
			return entities.DetailedProfile{}, fmt.Errorf("error while getting a profile from db: %w", err)
		}
		follows, err := getDBFollows(id)
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("while getting follows from db: %w", err)
		}
		detailedProfile := entities.DetailedProfile{Profile: profile, FollowsProfiles: follows}
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

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store_contracts.StoreProfileCreator {
	return createDBProfile
}

func NewStoreProfileGetter(getDBProfile DBProfileGetter) store_contracts.StoreProfileGetter {
	return getDBProfile
}

func NewStoreFollowsGetter(dbFollowsGetter DBFollowsGetter) store_contracts.StoreFollowsGetter {
	return dbFollowsGetter
}

func NewStoreFollowChecker(dbFollowChecker DBFollowChecker) store_contracts.StoreFollowChecker {
	return dbFollowChecker
}

func NewStoreFollower(dbFollower DBFollower) store_contracts.StoreFollower {
	return dbFollower
}

func NewStoreUnfollower(dbUnfollower DBUnfollower) store_contracts.StoreUnfollower {
	return dbUnfollower
}
