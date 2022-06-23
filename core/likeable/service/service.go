package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_values"
)

type (
	StoreLikeChecker          func(targetId string, fromUser core_values.UserId) (bool, error)
	StoreLike                 func(targetId string, fromUser core_values.UserId) error
	StoreUnlike               func(targetId string, fromUser core_values.UserId) error
	StoreLikesCountGetter     func(targetId string) (int, error)
	StoreUserLikesCountGetter func(id core_values.UserId) (int, error)
	StoreUserLikesGetter      func(id core_values.UserId) ([]string, error)
)

type (
	LikeToggler          func(target string, owner, liker core_values.UserId) error
	LikesCountGetter     func(targetId string) (int, error)
	UserLikesCountGetter func(core_values.UserId) (int, error)
	UserLikesGetter      func(core_values.UserId) ([]string, error)
	LikeChecker          func(targetId string, fromUser core_values.UserId) (bool, error)
)

func NewLikeToggler(checkLiked StoreLikeChecker, like StoreLike, unlike StoreUnlike) LikeToggler {
	return func(target string, owner, fromUser core_values.UserId) error {
		if owner == fromUser {
			return client_errors.LikingYourself
		}
		isLiked, err := checkLiked(target, fromUser)
		if err != nil {
			return fmt.Errorf("while checking if the target Likeable is liked: %w", err)
		}

		if isLiked {
			err = unlike(target, fromUser)
			if err != nil {
				return fmt.Errorf("while unliking a Likeable in service: %w", err)
			}
		} else {
			err = like(target, fromUser)
			if err != nil {
				return fmt.Errorf("while liking a Likeable in service: %w", err)
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
