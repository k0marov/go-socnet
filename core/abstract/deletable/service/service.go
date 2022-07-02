package service

import "github.com/k0marov/go-socnet/core/general/core_values"

type StoreDeleter func(targetId string, caller core_values.UserId) error

type Deleter func(targetId string, caller core_values.UserId) error

func NewDeleter(delete StoreDeleter) Deleter {
	return Deleter(delete)
}
