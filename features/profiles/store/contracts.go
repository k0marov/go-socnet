package store

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/core_values/ref"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"
)

type (
	AvatarFileCreator func(data ref.Ref[[]byte], belongsToUser core_values.UserId) (string, error)

	DBProfileGetter  func(id core_values.UserId) (models.ProfileModel, error)
	DBProfileCreator func(models.ProfileModel) error
	DBProfileUpdater func(id core_values.UserId, updData DBUpdateData) error

	DBFollowsGetter func(id core_values.UserId) ([]core_values.UserId, error)
	DBFollowChecker func(target, follower core_values.UserId) (bool, error)
	DBFollower      func(target, follower core_values.UserId) error
	DBUnfollower    func(target, unfollower core_values.UserId) error
)
