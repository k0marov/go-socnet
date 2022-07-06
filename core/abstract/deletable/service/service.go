package service

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type StoreDeleter func(targetId string) error

type Deleter func(targetId string, caller core_values.UserId) error
type ForceDeleter func(targetId string) error

func NewDeleter(getOwner ownable.OwnerGetter, delete StoreDeleter) Deleter {
	return func(targetId string, caller core_values.UserId) error {
		owner, err := getOwner(targetId)
		if err != nil {
			return core_err.Rethrow("getting owner of target", err)
		}
		if caller != owner {
			return client_errors.InsufficientPermissions
		}
		err = delete(targetId)
		if err != nil {
			return core_err.Rethrow("deleting the target", err)
		}
		return nil
	}
}

func NewForceDeleter(delete StoreDeleter) ForceDeleter {
	return ForceDeleter(delete)
}
