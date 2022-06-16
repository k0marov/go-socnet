package service

import (
	"core/client_errors"
	"core/core_errors"
	"core/core_values"
	core_entities "core/entities"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/store"
	"profiles/domain/validators"
	"profiles/domain/values"
)

type (
	ProfileGetter  func(id, caller core_values.UserId) (entities.ContextedProfile, error)
	FollowsGetter  func(core_values.UserId) ([]entities.Profile, error)
	FollowToggler  func(target, follower core_values.UserId) error
	ProfileUpdater func(core_entities.User, values.ProfileUpdateData) (entities.Profile, error)
	AvatarUpdater  func(core_entities.User, values.AvatarData) (values.AvatarPath, error)
	ProfileCreator func(core_entities.User) (entities.Profile, error)
)

func NewProfileGetter(getProfile store.StoreProfileGetter, isFollowed store.StoreFollowChecker) ProfileGetter {
	return func(id core_values.UserId, caller core_values.UserId) (entities.ContextedProfile, error) {
		profile, err := getProfile(id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.ContextedProfile{}, client_errors.NotFound
			}
			return entities.ContextedProfile{}, fmt.Errorf("while getting profile in a service: %w", err)
		}
		followedByCaller, err := isFollowed(id, caller)
		if err != nil {
			return entities.ContextedProfile{}, fmt.Errorf("while checking if profile is followed by caller: %w", err)
		}
		contextedProfile := entities.ContextedProfile{Profile: profile, IsFollowedByCaller: followedByCaller}
		return contextedProfile, nil
	}
}

func NewFollowsGetter(storeFollowsGetter store.StoreFollowsGetter) FollowsGetter {
	return func(userId core_values.UserId) ([]entities.Profile, error) {
		follows, err := storeFollowsGetter(userId)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return []entities.Profile{}, client_errors.NotFound
			}
			return []entities.Profile{}, fmt.Errorf("while getting follows in service: %w", err)
		}
		return follows, nil
	}
}

func NewFollowToggler(storeFollowChecker store.StoreFollowChecker, storeFollower store.StoreFollower, storeUnfollower store.StoreUnfollower) FollowToggler {
	return func(target, follower core_values.UserId) error {
		if target == follower {
			return client_errors.FollowingYourself
		}
		isFollowed, err := storeFollowChecker(target, follower)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return client_errors.NotFound
			}
			return fmt.Errorf("while checking if target is already followed: %w", err)
		}
		if isFollowed {
			err = storeUnfollower(target, follower)
			if err != nil {
				return fmt.Errorf("while unfollowing target: %w", err)
			}
		} else {
			err = storeFollower(target, follower)
			if err != nil {
				return fmt.Errorf("while following target: %w", err)
			}
		}
		return nil
	}
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

// this should be invoked when a new user is registered
func NewProfileCreator(storeProfileCreator store.StoreProfileCreator) ProfileCreator {
	return func(user core_entities.User) (entities.Profile, error) {
		newProfile := values.NewProfile{
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
			Id:         newProfile.Id,
			Username:   newProfile.Username,
			About:      newProfile.About,
			AvatarPath: newProfile.AvatarPath,
			Follows:    0,
			Followers:  0,
		}
		return createdProfile, nil
	}
}

func NewAvatarUpdater(validator validators.AvatarValidator, storeAvatar store.StoreAvatarUpdater) AvatarUpdater {
	return func(user core_entities.User, avatar values.AvatarData) (values.AvatarPath, error) {
		if clientError, ok := validator(avatar); !ok {
			return values.AvatarPath{}, clientError
		}

		avatarURL, err := storeAvatar(user.Id, avatar)
		if err != nil {
			return values.AvatarPath{}, fmt.Errorf("got an error while storing updated avatar: %w", err)
		}

		return avatarURL, nil
	}
}
