package store_contracts

import (
	"profiles/domain/entities"
	"profiles/domain/values"
)

type StoreProfileUpdater = func(id string, upd values.ProfileUpdateData) (entities.DetailedProfile, error)
type StoreDetailedProfileGetter = func(id string) (entities.DetailedProfile, error)
type StoreProfileCreator = func(values.NewProfile) error
type StoreAvatarUpdater = func(userId string, avatar values.AvatarData) (values.AvatarURL, error)