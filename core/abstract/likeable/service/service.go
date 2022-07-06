package service

import (
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

// TODO: add checks for core_err.ErrNotFound to all services

type (
	StoreLikeChecker          func(targetId string, fromUser core_values.UserId) (bool, error)
	StoreLike                 func(targetId string, fromUser core_values.UserId) error
	StoreUnlike               func(targetId string, fromUser core_values.UserId) error
	StoreLikesCountGetter     func(targetId string) (int, error)
	StoreUserLikesCountGetter func(id core_values.UserId) (int, error)
	StoreUserLikesGetter      func(id core_values.UserId) ([]string, error)
)

type (
	LikeToggler          func(target string, liker core_values.UserId) error
	LikesCountGetter     func(targetId string) (int, error)
	UserLikesCountGetter func(core_values.UserId) (int, error)
	UserLikesGetter      func(core_values.UserId) ([]string, error)
	LikeChecker          func(targetId string, fromUser core_values.UserId) (bool, error)
)

func NewLikeToggler(checkLiked StoreLikeChecker, like StoreLike, unlike StoreUnlike) LikeToggler {
	return func(target string, fromUser core_values.UserId) error {
		isLiked, err := checkLiked(target, fromUser)
		if err != nil {
			return core_err.Rethrow("checking if the target Likeable is liked", err)
		}

		if isLiked {
			err = unlike(target, fromUser)
			if err != nil {
				return core_err.Rethrow("unliking a Likeable in service", err)
			}
		} else {
			err = like(target, fromUser)
			if err != nil {
				return core_err.Rethrow("liking a Likeable in service", err)
			}
		}

		return nil
	}
}

func NewLikesCountGetter(getLikesCount StoreLikesCountGetter) LikesCountGetter {
	return LikesCountGetter(getLikesCount)
}

func NewUserLikesCountGetter(getUserLikesCount StoreUserLikesCountGetter) UserLikesCountGetter {
	return UserLikesCountGetter(getUserLikesCount)
}
func NewUserLikesGetter(getUserLikes StoreUserLikesGetter) UserLikesGetter {
	return UserLikesGetter(getUserLikes)
}

func NewLikeChecker(checkLiked StoreLikeChecker) LikeChecker {
	return LikeChecker(checkLiked)
}
