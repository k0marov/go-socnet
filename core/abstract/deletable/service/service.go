package service

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type StoreDeleter func(targetId string) error

type Deleter func(targetId string, caller core_values.UserId) error

func NewDeleter(getOwner ownable.OwnerGetter, delete StoreDeleter) Deleter {
	return func(targetId string, caller core_values.UserId) error {
		panic("unimplemented")
	}
}
