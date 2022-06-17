package store

import (
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"

	"github.com/k0marov/socnet/core/core_values"
)

type (
	StoreFollowChecker         func(target, follower core_values.UserId) (bool, error)
	StoreFollower              func(target, follower core_values.UserId) error
	StoreUnfollower            func(target, unfollower core_values.UserId) error
	StoreFollowsGetter         func(core_values.UserId) ([]core_values.UserId, error)
	StoreProfileGetter         func(id core_values.UserId) (entities.Profile, error)
	StoreProfileUpdater        func(id core_values.UserId, upd values.ProfileUpdateData) (entities.Profile, error)
	StoreDetailedProfileGetter func(id core_values.UserId) (entities.Profile, error)
	StoreProfileCreator        func(values.NewProfile) error
	StoreAvatarUpdater         func(userId core_values.UserId, avatar values.AvatarData) (core_values.FileURL, error)
)
