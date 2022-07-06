package service

import (
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type StoreOwnerGetter func(targetId string) (core_values.UserId, error)

type OwnerGetter func(targetId string) (core_values.UserId, error)

func NewOwnerGetter(getOwner StoreOwnerGetter) OwnerGetter {
	return func(target string) (core_values.UserId, error) {
		owner, err := getOwner(target)
		if err == core_err.ErrNotFound {
			return "", client_errors.NotFound
		}
		if err != nil {
			return "", core_err.Rethrow("getting owner from store", err)
		}
		return owner, nil
	}
}
