package service

import "github.com/k0marov/go-socnet/core/general/core_values"

type StoreRecsGetter = func(user core_values.UserId, count int) ([]string, error)
type StoreRecsSetter = func(core_values.UserId, []string) error

type RecsGetter = func(user core_values.UserId, count int) ([]string, error)
type RecsUpdater = func() error

func NewRecsUpdater() RecsUpdater {
	return func() error {
		panic("unimplemented")
	}
}

func NewRecsGetter(getRecs StoreRecsGetter) RecsGetter {
	// TODO: complete with random posts if store returns not enough
	return RecsGetter(getRecs)
}
