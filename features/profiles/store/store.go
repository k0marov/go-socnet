package store

import (
	"core/core_values"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store"
	"profiles/domain/values"
)

type DBUpdateData struct {
	About      string
	AvatarPath string
}

func NewStoreAvatarUpdater(createFile AvatarFileCreator, updateDBProfile DBProfileUpdater) store.StoreAvatarUpdater {
	return func(userId core_values.UserId, avatar values.AvatarData) (values.AvatarPath, error) {
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

func NewStoreProfileUpdater(updateDBProfile DBProfileUpdater, getProfile store.StoreDetailedProfileGetter) store.StoreProfileUpdater {
	return func(id core_values.UserId, upd values.ProfileUpdateData) (entities.Profile, error) {
		err := updateDBProfile(id, DBUpdateData{About: upd.About})
		if err != nil {
			return entities.Profile{}, fmt.Errorf("error while updating profile in db: %w", err)
		}
		return getProfile(id)
	}
}

func NewStoreProfileCreator(createDBProfile DBProfileCreator) store.StoreProfileCreator {
	return store.StoreProfileCreator(createDBProfile)
}

func NewStoreProfileGetter(getDBProfile DBProfileGetter) store.StoreProfileGetter {
	return store.StoreProfileGetter(getDBProfile)
}

func NewStoreFollowsGetter(dbFollowsGetter DBFollowsGetter) store.StoreFollowsGetter {
	return store.StoreFollowsGetter(dbFollowsGetter)
}

func NewStoreFollowChecker(dbFollowChecker DBFollowChecker) store.StoreFollowChecker {
	return store.StoreFollowChecker(dbFollowChecker)
}

func NewStoreFollower(dbFollower DBFollower) store.StoreFollower {
	return store.StoreFollower(dbFollower)
}

func NewStoreUnfollower(dbUnfollower DBUnfollower) store.StoreUnfollower {
	return store.StoreUnfollower(dbUnfollower)
}
