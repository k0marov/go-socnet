package service

import (
	"github.com/k0marov/go-socnet/core/abstract/recommendable"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"strconv"
)

type FeedGetter = func(count string, caller core_values.UserId) ([]string, error)

const DefaultCount = 5
const MaxCount = 50

func convertCount(countStr string) (count int, ok bool) {
	if countStr == "" {
		return DefaultCount, true
	}
	countConv, err := strconv.Atoi(countStr)
	return countConv, err == nil
}

func NewFeedGetter(getFeed recommendable.RecsGetter) FeedGetter {
	return func(countStr string, caller core_values.UserId) ([]string, error) {
		count, ok := convertCount(countStr)
		if !ok {
			return []string{}, client_errors.NonIntegerCount
		}
		if count > MaxCount {
			return []string{}, client_errors.TooBigCount
		}
		return getFeed(caller, count)
	}
}
