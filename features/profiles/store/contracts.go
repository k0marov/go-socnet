package store

import (
	"core/ref"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type AvatarFileCreator = func(data ref.Ref[[]byte], belongsToUser values.UserId) (string, error)

type (
	DBProfileGetter  = func(id values.UserId) (entities.Profile, error)
	DBFollowsGetter  = func(id values.UserId) ([]entities.Profile, error)
	DBFollowChecker  = func(target, follower values.UserId) (bool, error)
	DBFollower       = func(target, follower values.UserId) error
	DBUnfollower     = func(target, unfollower values.UserId) error
	DBProfileCreator = func(entities.Profile) error
	DBProfileUpdater = func(id values.UserId, updData DBUpdateData) error
)
