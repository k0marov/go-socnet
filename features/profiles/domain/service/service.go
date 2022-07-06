package service

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_entities"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/static_store"

	"github.com/k0marov/go-socnet/features/profiles/domain/contexters"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
	"github.com/k0marov/go-socnet/features/profiles/domain/store"
	"github.com/k0marov/go-socnet/features/profiles/domain/validators"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"
)

type (
	ProfileGetter  func(id, caller core_values.UserId) (entities.ContextedProfile, error)
	ProfileUpdater func(core_entities.User, values.ProfileUpdateData) (entities.ContextedProfile, error)
	AvatarUpdater  func(core_entities.User, values.AvatarData) (core_values.FileURL, error)
	ProfileCreator func(core_entities.User) (entities.Profile, error)
	FollowToggler  func(target, follower core_values.UserId) error
	FollowsGetter  func(target, caller core_values.UserId) ([]entities.ContextedProfile, error)
)

func NewProfileGetter(getProfile store.StoreProfileGetter, addContext contexters.ProfileContextAdder) ProfileGetter {
	return func(id core_values.UserId, caller core_values.UserId) (entities.ContextedProfile, error) {
		profile, err := getProfile(id)
		if err != nil {
			if err == core_err.ErrNotFound {
				return entities.ContextedProfile{}, client_errors.NotFound
			}
			return entities.ContextedProfile{}, core_err.Rethrow("getting profile in a service", err)
		}
		contextedProfile, err := addContext(profile, caller)
		if err != nil {
			return entities.ContextedProfile{}, core_err.Rethrow("adding context to profile", err)
		}
		return contextedProfile, nil
	}
}

func NewFollowToggler(toggleLike likeable.LikeToggler) FollowToggler {
	return FollowToggler(toggleLike)
}

func NewFollowsGetter(getUserLikes likeable.UserLikesGetter, getProfile ProfileGetter) FollowsGetter {
	return func(target, caller core_values.UserId) ([]entities.ContextedProfile, error) {
		followIds, err := getUserLikes(target)
		if err != nil {
			return []entities.ContextedProfile{}, core_err.Rethrow("getting a list of profile ids that target follows", err)
		}
		var follows []entities.ContextedProfile
		for _, followId := range followIds {
			follow, err := getProfile(followId, caller)
			if err != nil {
				return []entities.ContextedProfile{}, fmt.Errorf("while getting profile details for a follow %s: %w", followId, err)
			}
			follows = append(follows, follow)
		}
		return follows, nil
	}
}

func NewProfileUpdater(validate validators.ProfileUpdateValidator, update store.StoreProfileUpdater, get ProfileGetter) ProfileUpdater {
	return func(user core_entities.User, updateData values.ProfileUpdateData) (entities.ContextedProfile, error) {
		if clientError, ok := validate(updateData); !ok {
			return entities.ContextedProfile{}, clientError
		}
		err := update(user.Id, updateData)
		if err != nil {
			return entities.ContextedProfile{}, fmt.Errorf("got an error while updating profile in a service: %w", err)
		}
		updatedProfile, err := get(user.Id, user.Id)
		if err != nil {
			return entities.ContextedProfile{}, err
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
