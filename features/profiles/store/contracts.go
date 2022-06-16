package store

import (
	"core/core_values"
	"core/ref"
	"profiles/domain/entities"
	"profiles/domain/values"
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
