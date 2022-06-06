package service_contracts

import (
	core_entities "core/entities"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type DetailedProfileGetter = func(core_entities.User) (entities.DetailedProfile, error)
type ProfileUpdater = func(core_entities.User, values.ProfileUpdateData) (entities.DetailedProfile, error)
type AvatarUpdater = func(core_entities.User, values.AvatarData) (values.AvatarURL, error)
type ProfileCreator = func(core_entities.User) (entities.DetailedProfile, error)
