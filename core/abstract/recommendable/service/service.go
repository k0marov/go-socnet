package service

import "github.com/k0marov/go-socnet/core/general/core_values"

type StoreRecsGetter = func(core_values.UserId) ([]string, error)

type RecsGetter = func(core_values.UserId) ([]string, error)

func NewRecsGetter(getRecs StoreRecsGetter) RecsGetter {
	return RecsGetter(getRecs)
}
