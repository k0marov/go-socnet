package store

import (
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"
	"github.com/k0marov/socnet/features/profiles/store/models"

	"github.com/k0marov/socnet/core/core_values"
)

type (
	StoreProfileGetter  func(id core_values.UserId) (entities.Profile, error)
	StoreProfileUpdater func(id core_values.UserId, upd values.ProfileUpdateData) (entities.Profile, error)
	StoreProfileCreator func(model models.ProfileModel) error
	StoreAvatarUpdater  func(userId core_values.UserId, avatar values.AvatarData) (core_values.FileURL, error)
)
