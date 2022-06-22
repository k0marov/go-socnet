package service

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable"
)

func NewLikeToggler() likeable.LikeToggler {
	return func(id string, fromUser core_values.UserId) error {
		panic("unimplemented")
	}
}

func NewLikesCountGetter() likeable.LikesCountGetter {
	return func(id string) (int, error) {
		panic("unimplemented")
	}
}

func NewLikeChecker() likeable.LikeChecker {
	return func(id string, fromUser core_values.UserId) (bool, error) {
		panic("unimplemented")
	}
}
