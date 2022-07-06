package service

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SafeLikeToggler func(target string, caller core_values.UserId) error

func NewSafeLikeToggler(getOwner ownable.OwnerGetter, toggleLike likeable.LikeToggler) SafeLikeToggler {
	return func(target string, caller core_values.UserId) error {
		owner, err := getOwner(target)
		if err != nil {
			return core_err.Rethrow("getting owner of OwnableLikeable", err)
		}
		if owner == caller {
			return client_errors.LikingYourself
		}
		err = toggleLike(target, caller)
		if err != nil {
			return core_err.Rethrow("toggling like on OwnableLikeable", err)
		}
		return nil
	}
}
