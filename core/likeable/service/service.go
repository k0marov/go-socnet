package service

import (
	"fmt"
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_values"
)

type (
	StoreLikeChecker      func(id string, fromUser core_values.UserId) (bool, error)
	StoreLike             func(id string, fromUser core_values.UserId) error
	StoreUnlike           func(id string, fromUser core_values.UserId) error
	StoreLikesCountGetter func(id string) (int, error)
)

type (
	LikeToggler      func(id string, owner, liker core_values.UserId) error
	LikesCountGetter func(id string) (int, error)
	LikeChecker      func(id string, fromUser core_values.UserId) (bool, error)
)

func NewLikeToggler(checkLiked StoreLikeChecker, like StoreLike, unlike StoreUnlike) LikeToggler {
	return func(id string, owner, fromUser core_values.UserId) error {
		if owner == fromUser {
			return client_errors.LikingYourself
		}
		isLiked, err := checkLiked(id, fromUser)
		if err != nil {
			return fmt.Errorf("while checking if the target Likeable is liked: %w", err)
		}

		if isLiked {
			err = unlike(id, fromUser)
			if err != nil {
				return fmt.Errorf("while unliking a Likable in service: %w", err)
			}
		} else {
			err = like(id, fromUser)
			if err != nil {
				return fmt.Errorf("while liking a Likeable in service: %w", err)
			}
		}

		return nil
	}
}

func NewLikesCountGetter(getLikesCount StoreLikesCountGetter) LikesCountGetter {
	return func(id string) (int, error) {
		likes, err := getLikesCount(id)
		if err != nil {
			return 0, fmt.Errorf("while getting count of Likeable likes in service: %w", err)
		}
		return likes, nil
	}
}

func NewLikeChecker(checkLiked StoreLikeChecker) LikeChecker {
	return func(id string, fromUser core_values.UserId) (bool, error) {
		isLiked, err := checkLiked(id, fromUser)
		if err != nil {
			return false, fmt.Errorf("while checking if Likeable is liked in service: %w", err)
		}
		return isLiked, nil
	}
}
