package store_contract

import (
	"errors"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type ProfileStore interface {
	GetByIdDetailed(userId string) (entities.DetailedProfile, error)
	StoreNew(entities.DetailedProfile) error
	Update(userId string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error)
	StoreAvatar(userId string, avatar values.AvatarData) (values.AvatarURL, error)
}

var ErrProfileNotFound = errors.New("profile not found")
