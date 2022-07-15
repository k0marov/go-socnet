package service

import (
	"errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type StoreRandomGetter = func(count int) ([]string, error)
type StoreRecsGetter = func(user core_values.UserId, count int) ([]string, error)
type StoreRecsSetter = func(core_values.UserId, []string) error

type RecsGetter = func(user core_values.UserId, count int) ([]string, error)
type RecsUpdater = func() error

func NewRecsUpdater() RecsUpdater {
	return func() error {
		return errors.New("unimplemented")
	}
}

func NewRecsGetter(getRecs StoreRecsGetter, getRandom StoreRandomGetter) RecsGetter {
	return func(user core_values.UserId, count int) ([]string, error) {
		recs, err := getRecs(user, count)
		if err != nil {
			return []string{}, core_err.Rethrow("getting recommendations", err)
		}
		if len(recs) == count {
			return recs, nil
		}

		randomRecs, err := getRandom(count - len(recs))
		if err != nil {
			return []string{}, core_err.Rethrow("getting random recommendations", err)
		}
		return append(recs, randomRecs...), nil
	}
}
