package service

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SafeLikeToggler func(target string, caller core_values.UserId) error

func NewSafeLikeToggler(getOwner ownable.OwnerGetter, toggleLike likeable.LikeToggler) SafeLikeToggler {
	return func(target string, caller core_values.UserId) error {
		owner, err := getOwner(target)
		if err != nil {
			return fmt.Errorf("while getting owner of OwnableLikeable: %w", err)
		}
		if owner == caller {
			return client_errors.LikingYourself
		}
		err = toggleLike(target, caller)
		if err != nil {
			return fmt.Errorf("while toggling like on OwnableLikeable: %w", err)
		}
		return nil
	}
}
