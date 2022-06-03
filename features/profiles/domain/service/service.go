package service

import (
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

func (p *ProfileService) GetOrCreateDetailed(user core_entities.User) (entities.DetailedProfile, error) {
	profile, err := p.store.GetByIdDetailed(user.Id)
	if err != nil {
		if err == ErrProfileNotFound {
			return p.createAndReturn(user)
		}
		return entities.DetailedProfile{}, fmt.Errorf("got an error while getting profile in a service: %w", err)
	}

	return profile, nil
}

const DefaultAbout = ""

func (p *ProfileService) createAndReturn(user core_entities.User) (entities.DetailedProfile, error) {
	newProfile := entities.DetailedProfile{
		Profile: entities.Profile{
			Id:       user.Id,
			Username: user.Username,
			About:    DefaultAbout,
		},
	}
	err := p.store.StoreNew(newProfile)
	if err != nil {
		return entities.DetailedProfile{}, fmt.Errorf("got an error while creating a profile in a service: %w", err)
	}
	return newProfile, nil
}
