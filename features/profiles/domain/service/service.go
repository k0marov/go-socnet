package service

import (
	"core/client_errors"
	"core/core_errors"
	"core/core_values"
	core_entities "core/entities"
	"fmt"
	"profiles/domain/entities"
	store "profiles/domain/store_contracts"
	"profiles/domain/validators"
	"profiles/domain/values"
)

type (
	ProfileGetter         = func(core_values.UserId) (entities.Profile, error)
	FollowsGetter         = func(core_values.UserId) ([]entities.Profile, error)
	FollowToggler         = func(target, follower core_values.UserId) error
	DetailedProfileGetter = func(core_entities.User) (entities.DetailedProfile, error)
	ProfileUpdater        = func(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
	AvatarUpdater         = func(core_entities.User, values.AvatarData) (values.AvatarPath, error)
	ProfileCreator        = func(core_entities.User) (entities.DetailedProfile, error)
)

func NewProfileGetter(storeProfileGetter store.StoreProfileGetter) ProfileGetter {
	return func(id core_values.UserId) (entities.Profile, error) {
		profile, err := storeProfileGetter(id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.Profile{}, client_errors.ProfileNotFound
			}
			return entities.Profile{}, fmt.Errorf("while getting profile in a service: %w", err)
		}
		return profile, nil
	}
}

func NewFollowsGetter(storeFollowsGetter store.StoreFollowsGetter) FollowsGetter {
	return func(userId core_values.UserId) ([]entities.Profile, error) {
		follows, err := storeFollowsGetter(userId)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return []entities.Profile{}, client_errors.ProfileNotFound
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
				return client_errors.ProfileNotFound
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
	return func(user core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
		if clientError, ok := validator(updateData); !ok {
			return entities.DetailedProfile{}, clientError
		}
		updatedProfile, err := storeProfileUpdater(user.Id, updateData)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.DetailedProfile{}, client_errors.ProfileNotFound
			}
			return entities.DetailedProfile{}, fmt.Errorf("got an error while updating profile in a service: %w", err)
		}
		return updatedProfile, nil
	}
}

func NewDetailedProfileGetter(storeDetailedGetter store.StoreDetailedProfileGetter) DetailedProfileGetter {
	return func(user core_entities.User) (entities.DetailedProfile, error) {
		profile, err := storeDetailedGetter(user.Id)
		if err != nil {
			if err == core_errors.ErrNotFound {
				return entities.DetailedProfile{}, client_errors.ProfileNotFound
			}
			return entities.DetailedProfile{}, fmt.Errorf("got an error while getting profile in a service: %w", err)
		}

		return profile, nil
	}
}

const DefaultAbout = ""
const DefaultAvatarPath = ""

// this should be invoked when a new user is registered
func NewProfileCreator(storeProfileCreator store.StoreProfileCreator) ProfileCreator {
	return func(user core_entities.User) (entities.DetailedProfile, error) {
		newProfile := values.NewProfile{
			Id:         user.Id,
			Username:   user.Username,
			About:      DefaultAbout,
			AvatarPath: DefaultAvatarPath,
		}
		err := storeProfileCreator(newProfile)
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("got an error while creating a profile in a service: %w", err)
		}
		createdProfile := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         newProfile.Id,
				Username:   newProfile.Username,
				About:      newProfile.About,
				AvatarPath: newProfile.AvatarPath,
				Follows:    0,
				Followers:  0,
			}, FollowsProfiles: []entities.Profile{},
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
