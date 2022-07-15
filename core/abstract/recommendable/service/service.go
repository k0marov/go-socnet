package service

import "github.com/k0marov/go-socnet/core/general/core_values"

type StoreRecsGetter = func(user core_values.UserId, count int) ([]string, error)
type StoreRecsSetter = func(core_values.UserId, []string) error

type RecsGetter = func(user core_values.UserId, count int) ([]string, error)
type RecsUpdater = func(core_values.UserId) error

func NewRecsUpdater() RecsUpdater {
	return func(id core_values.UserId) error {
		panic("unimplemented")
	}
}

func NewRecsGetter(getRecs StoreRecsGetter) RecsGetter {
	return RecsGetter(getRecs)
}
