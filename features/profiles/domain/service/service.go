package service

import (
	"core/client_errors"
	"core/core_errors"
	core_entities "core/entities"
	"core/image_decoder"
	"fmt"
	"profiles/domain/entities"
	contracts "profiles/domain/service_contracts"
	"profiles/domain/values"
)

type StoreProfileUpdater = func(id string, upd values.ProfileUpdateData) (entities.DetailedProfile, error)

const MaxAboutLength = 255

func NewProfileUpdater(storeProfileUpdater StoreProfileUpdater) contracts.ProfileUpdater {
	return func(user core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
		if len(updateData.About) > MaxAboutLength {
			return entities.DetailedProfile{}, client_errors.AboutTooLong
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

type StoreDetailedProfileGetter = func(id string) (entities.DetailedProfile, error)

func NewDetailedProfileGetter(storeDetailedGetter StoreDetailedProfileGetter) contracts.DetailedProfileGetter {
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

type StoreProfileCreator = func(entities.DetailedProfile) error

const DefaultAbout = ""
const DefaultAvatarPath = ""

// this should be invoked when a new user is registered
func NewProfileCreator(storeProfileCreator StoreProfileCreator) contracts.ProfileCreator {
	return func(user core_entities.User) (entities.DetailedProfile, error) {
		newProfile := entities.DetailedProfile{
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
		return newProfile, nil
	}
}

type StoreAvatarUpdater = func(userId string, avatar values.AvatarData) (values.AvatarURL, error)

func NewAvatarUpdater(storeAvatarUpdater StoreAvatarUpdater, imageDecoder image_decoder.ImageDecoder) contracts.AvatarUpdater {
	return func(user core_entities.User, avatar values.AvatarData) (values.AvatarURL, error) {
		imageDimensions, err := imageDecoder(avatar.Data)
		if err != nil {
			return values.AvatarURL{}, client_errors.NonImageAvatar
		}
		if imageDimensions.Height != imageDimensions.Width {
			return values.AvatarURL{}, client_errors.NonSquareAvatar
		}

		avatarURL, err := storeAvatarUpdater(user.Id, avatar)
		if err != nil {
			return values.AvatarURL{}, fmt.Errorf("got an error while storing updated avatar: %w", err)
		}

		return avatarURL, nil
	}
}
