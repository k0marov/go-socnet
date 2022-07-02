package service

import "github.com/k0marov/go-socnet/core/general/core_values"

type StoreOwnerGetter func(targetId string) (core_values.UserId, error)

type OwnerGetter func(targetId string) (core_values.UserId, error)

func NewOwnerGetter(getOwner StoreOwnerGetter) OwnerGetter {
	return func(targetId string) (core_values.UserId, error) {
		panic("unimplemented")
	}
}
