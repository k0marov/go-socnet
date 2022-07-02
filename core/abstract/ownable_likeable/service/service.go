package service

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SafeLikeToggler func(target string, caller core_values.UserId) error

func NewSafeLikeToggler(getOwner ownable.OwnerGetter, toggleLike likeable.LikeToggler) SafeLikeToggler {
	return func(target string, caller core_values.UserId) error {
		return nil
	}
}
