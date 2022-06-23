package store

import (
	"fmt"
	"github.com/k0marov/socnet/core/likeable"

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
	return func(userId core_values.UserId, avatar values.AvatarData) (core_values.StaticPath, error) {
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

func NewStoreProfileUpdater(updateDBProfile DBProfileUpdater, getProfile store.StoreProfileGetter) store.StoreProfileUpdater {
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

func NewStoreProfileGetter(getDBProfile DBProfileGetter, getFollowers likeable.LikesCountGetter, getFollows likeable.UserLikesCountGetter) store.StoreProfileGetter {
	return func(id core_values.UserId) (entities.Profile, error) {
		profileModel, err := getDBProfile(id)
		if err != nil {
			return entities.Profile{}, fmt.Errorf("while getting the profile model from db: %w", err)
		}
		followers, err := getFollowers(id)
		if err != nil {
			return entities.Profile{}, err
		}
		follows, err := getFollows(id)
		if err != nil {
			return entities.Profile{}, err
		}
		profile := entities.Profile{
			ProfileModel: profileModel,
			Follows:      follows,
			Followers:    followers,
		}
		return profile, nil
	}
}
