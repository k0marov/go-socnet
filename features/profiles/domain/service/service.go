package service

import (
	"core/client_errors"
	"core/core_errors"
	core_entities "core/entities"
	"fmt"
	"profiles/domain/entities"
	store "profiles/domain/store_contracts"
	"profiles/domain/validators"
	"profiles/domain/values"
)

type DetailedProfileGetter = func(core_entities.User) (entities.DetailedProfile, error)
type ProfileUpdater = func(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
type AvatarUpdater = func(core_entities.User, values.AvatarData) (values.AvatarURL, error)
type ProfileCreator = func(core_entities.User) (entities.DetailedProfile, error)

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
			Profile: entities.Profile{
				Id:         user.Id,
				Username:   user.Username,
				About:      DefaultAbout,
				AvatarPath: DefaultAvatarPath,
			},
		}
		err := storeProfileCreator(newProfile)
		if err != nil {
			return entities.DetailedProfile{}, fmt.Errorf("got an error while creating a profile in a service: %w", err)
		}
		createdProfile := entities.DetailedProfile{
			Profile: newProfile.Profile,
		}
		return createdProfile, nil
	}
}

func NewAvatarUpdater(validator validators.AvatarValidator, storeAvatar store.StoreAvatarUpdater) AvatarUpdater {
	return func(user core_entities.User, avatar values.AvatarData) (values.AvatarURL, error) {
		if clientError, ok := validator(avatar); !ok {
			return values.AvatarURL{}, clientError
		}

		avatarURL, err := storeAvatar(user.Id, avatar)
		if err != nil {
			return values.AvatarURL{}, fmt.Errorf("got an error while storing updated avatar: %w", err)
		}

		return avatarURL, nil
	}
}
