package store_contracts

import (
	"profiles/domain/entities"
	"profiles/domain/values"
)

type (
	StoreFollowChecker         = func(target, follower values.UserId) (bool, error)
	StoreFollower              = func(target, follower values.UserId) error
	StoreUnfollower            = func(target, unfollower values.UserId) error
	StoreFollowsGetter         = func(values.UserId) ([]entities.Profile, error)
	StoreProfileGetter         = func(id values.UserId) (entities.Profile, error)
	StoreProfileUpdater        = func(id values.UserId, upd values.ProfileUpdateData) (entities.DetailedProfile, error)
	StoreDetailedProfileGetter = func(id values.UserId) (entities.DetailedProfile, error)
	StoreProfileCreator        = func(entities.Profile) error
	StoreAvatarUpdater         = func(userId values.UserId, avatar values.AvatarData) (values.AvatarPath, error)
)
