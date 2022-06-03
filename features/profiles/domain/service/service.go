package service

import (
	"core/client_errors"
	core_entities "core/entities"
	"errors"
	"fmt"
	"profiles/domain/entities"
)

type ProfileStore interface {
	GetByIdDetailed(string) (entities.DetailedProfile, error)
	StoreNew(entities.DetailedProfile) error
}

var ErrProfileNotFound = errors.New("profile not found")

type ProfileService struct {
	store ProfileStore
}

func NewProfileService(store ProfileStore) *ProfileService {
	return &ProfileService{store}
}

func (p *ProfileService) GetDetailed(user core_entities.User) (entities.DetailedProfile, error) {
	profile, err := p.store.GetByIdDetailed(user.Id)
	if err != nil {
		if err == ErrProfileNotFound {
			return entities.DetailedProfile{}, client_errors.ProfileNotFound
			// return p.createAndReturn(user)
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
