package store

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/static_store"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
	"github.com/k0marov/go-socnet/features/profiles/domain/store"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"
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

func NewStoreProfileUpdater(updateDBProfile DBProfileUpdater) store.StoreProfileUpdater {
	return func(id core_values.UserId, upd values.ProfileUpdateData) error {
		err := updateDBProfile(id, DBUpdateData{About: upd.About})
		if err != nil {
			return fmt.Errorf("error while updating profile in db: %w", err)
		}
		return nil
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
			AvatarURL:    static_store.PathToURL(profileModel.AvatarPath),
			Follows:      follows,
			Followers:    followers,
		}
		return profile, nil
	}
}
