package service

import (
	"fmt"

	"github.com/k0marov/go-socnet/core/likeable"
	"github.com/k0marov/go-socnet/core/static_store"
	"github.com/k0marov/go-socnet/features/profiles/domain/contexters"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
	"github.com/k0marov/go-socnet/features/profiles/domain/store"
	"github.com/k0marov/go-socnet/features/profiles/domain/validators"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	core_entities "github.com/k0marov/go-socnet/core/core_entities"
	"github.com/k0marov/go-socnet/core/core_errors"
	"github.com/k0marov/go-socnet/core/core_values"
)

type (
	ProfileGetter  func(id, caller core_values.UserId) (entities.ContextedProfile, error)
	ProfileUpdater func(core_entities.User, values.ProfileUpdateData) (entities.Profile, error)
	AvatarUpdater  func(core_entities.User, values.AvatarData) (core_values.FileURL, error)
	ProfileCreator func(core_entities.User) (entities.Profile, error)
	FollowToggler  func(target, follower core_values.UserId) error
	FollowsGetter  func(target core_values.UserId) ([]core_values.UserId, error)
)

func NewProfileGetter(getProfile store.StoreProfileGetter, addContext contexters.ProfileContextAdder) ProfileGetter {
	return func(id core_values.UserId, caller core_values.UserId) (entities.ContextedProfile, error) {
		profile, err := getProfile(id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.ContextedProfile{}, client_errors.NotFound
			}
			return entities.ContextedProfile{}, fmt.Errorf("while getting profile in a service: %w", err)
		}
		contextedProfile, err := addContext(profile, caller)
		if err != nil {
			return entities.ContextedProfile{}, fmt.Errorf("while adding context to profile: %w", err)
		}
		return contextedProfile, nil
	}
}

func NewFollowToggler(toggleLike likeable.LikeToggler) FollowToggler {
	return func(target, follower core_values.UserId) error {
		return toggleLike(target, target, follower)
	}
}

func NewFollowsGetter(getUserLikes likeable.UserLikesGetter) FollowsGetter {
	return FollowsGetter(getUserLikes)
}

func NewProfileUpdater(validator validators.ProfileUpdateValidator, storeProfileUpdater store.StoreProfileUpdater) ProfileUpdater {
	return func(user core_entities.User, updateData values.ProfileUpdateData) (entities.Profile, error) {
		if clientError, ok := validator(updateData); !ok {
			return entities.Profile{}, clientError
		}
		updatedProfile, err := storeProfileUpdater(user.Id, updateData)
		if err != nil {
			return entities.Profile{}, fmt.Errorf("got an error while updating profile in a service: %w", err)
		}
		return updatedProfile, nil
	}
}

const DefaultAbout = ""
const DefaultAvatarPath = ""

// NewProfileCreator this should be invoked when a new user is registered
func NewProfileCreator(storeProfileCreator store.StoreProfileCreator) ProfileCreator {
	return func(user core_entities.User) (entities.Profile, error) {
		newProfile := models.ProfileModel{
			Id:         user.Id,
			Username:   user.Username,
			About:      DefaultAbout,
			AvatarPath: DefaultAvatarPath,
		}
		err := storeProfileCreator(newProfile)
		if err != nil {
			return entities.Profile{}, fmt.Errorf("got an error while creating a profile in a service: %w", err)
		}
		createdProfile := entities.Profile{
			ProfileModel: newProfile,
			Follows:      0,
			Followers:    0,
		}
		return createdProfile, nil
	}
}

func NewAvatarUpdater(validator validators.AvatarValidator, storeAvatar store.StoreAvatarUpdater) AvatarUpdater {
	return func(user core_entities.User, avatar values.AvatarData) (core_values.FileURL, error) {
		if clientError, ok := validator(avatar); !ok {
			return "", clientError
		}

		avatarPath, err := storeAvatar(user.Id, avatar)
		if err != nil {
			return "", fmt.Errorf("got an error while storing updated avatar: %w", err)
		}

		return static_store.PathToURL(avatarPath), nil
	}
}
