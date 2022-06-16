package store

import (
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"

	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/ref"
)

type (
	AvatarFileCreator func(data ref.Ref[[]byte], belongsToUser core_values.UserId) (string, error)

	DBProfileGetter  func(id core_values.UserId) (entities.Profile, error)
	DBProfileCreator func(values.NewProfile) error
	DBProfileUpdater func(id core_values.UserId, updData DBUpdateData) error

	DBFollowsGetter func(id core_values.UserId) ([]core_values.UserId, error)
	DBFollowChecker func(target, follower core_values.UserId) (bool, error)
	DBFollower      func(target, follower core_values.UserId) error
	DBUnfollower    func(target, unfollower core_values.UserId) error
)
