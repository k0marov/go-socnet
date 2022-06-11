package store_contracts

import (
	"profiles/domain/entities"
	"profiles/domain/values"
)

type StoreFollower = func(target, follower values.UserId) error
type StoreUnfollower = func(target, unfollower values.UserId) error
type StoreFollowsGetter = func(values.UserId) ([]entities.Profile, error)
type StoreProfileGetter = func(id values.UserId) (entities.Profile, error)
type StoreProfileUpdater = func(id values.UserId, upd values.ProfileUpdateData) (entities.DetailedProfile, error)
type StoreDetailedProfileGetter = func(id values.UserId) (entities.DetailedProfile, error)
type StoreProfileCreator = func(entities.Profile) error
type StoreAvatarUpdater = func(userId values.UserId, avatar values.AvatarData) (values.AvatarPath, error)
