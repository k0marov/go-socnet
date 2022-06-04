package service

import (
	"core/client_errors"
	core_entities "core/entities"
	"core/image_decoder"
	"errors"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type ProfileStore interface {
	GetByIdDetailed(userId string) (entities.DetailedProfile, error)
	StoreNew(entities.DetailedProfile) error
	Update(userId string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error)
	StoreAvatar(userId string, avatar values.AvatarData) (entities.DetailedProfile, error)
}

var ErrProfileNotFound = errors.New("profile not found")

type ProfileService struct {
	store        ProfileStore
	imageDecoder image_decoder.ImageDecoder
}

func NewProfileService(store ProfileStore, imageDecoder image_decoder.ImageDecoder) *ProfileService {
	return &ProfileService{store, imageDecoder}
}

const MaxAboutLength = 255

func (p *ProfileService) Update(user core_entities.User, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
	if len(updateData.About) > MaxAboutLength {
		return entities.DetailedProfile{}, client_errors.AboutTooLong
	}
	updatedProfile, err := p.store.Update(user.Id, updateData)
	if err != nil {
		if err == ErrProfileNotFound {
			return entities.DetailedProfile{}, client_errors.ProfileNotFound
		}
		return entities.DetailedProfile{}, fmt.Errorf("got an error while updating profile in a service: %w", err)
	}
	return updatedProfile, nil
}

func (p *ProfileService) UpdateAvatar(user core_entities.User, avatar values.AvatarData) (entities.DetailedProfile, error) {
	imageDimensions, err := p.imageDecoder.Decode(avatar.Data)
	if err != nil {
		return entities.DetailedProfile{}, client_errors.NonImageAvatar
	}
	if imageDimensions.Height != imageDimensions.Width {
		return entities.DetailedProfile{}, client_errors.NonSquareAvatar
	}

	updatedProfile, err := p.store.StoreAvatar(user.Id, avatar)
	if err != nil {
		return entities.DetailedProfile{}, fmt.Errorf("got an error while storing updated avatar: %w", err)
	}

	return updatedProfile, nil
}

func (p *ProfileService) GetDetailed(user core_entities.User) (entities.DetailedProfile, error) {
	profile, err := p.store.GetByIdDetailed(user.Id)
	if err != nil {
		if err == ErrProfileNotFound {
			return entities.DetailedProfile{}, client_errors.ProfileNotFound
		}
		return entities.DetailedProfile{}, fmt.Errorf("got an error while getting profile in a service: %w", err)
	}

	return profile, nil
}

const DefaultAbout = ""
const DefaultAvatarPath = ""

// this should be invoked when a new user is registered
func (p *ProfileService) CreateProfileForUser(user core_entities.User) (entities.DetailedProfile, error) {
	newProfile := entities.DetailedProfile{
		Profile: entities.Profile{
			Id:         user.Id,
			Username:   user.Username,
			About:      DefaultAbout,
			AvatarPath: DefaultAvatarPath,
		},
	}
	err := p.store.StoreNew(newProfile)
	if err != nil {
		return entities.DetailedProfile{}, fmt.Errorf("got an error while creating a profile in a service: %w", err)
	}
	return newProfile, nil
}
