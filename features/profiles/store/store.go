package store

import (
	"fmt"

	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/store"
	"github.com/k0marov/socnet/features/profiles/domain/values"

	"github.com/k0marov/socnet/core/core_values"
)

type DBUpdateData struct {
	About      string
	AvatarPath string
}

func NewStoreAvatarUpdater(createFile AvatarFileCreator, updateDBProfile DBProfileUpdater) store.StoreAvatarUpdater {
	return func(userId core_values.UserId, avatar values.AvatarData) (core_values.ImageUrl, error) {
		avatarUrl, err := createFile(avatar.Data, userId)
		if err != nil {
			return "", fmt.Errorf("error while storing avatar file: %w", err)
		}
		updateData := DBUpdateData{
			AvatarPath: avatarUrl,
		}
		err = updateDBProfile(userId, updateData)
		if err != nil {
			return "", fmt.Errorf("error while updating avatar path in DB: %w", err)
		}
		return avatarUrl, nil
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
