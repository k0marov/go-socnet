package service_test

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable/service"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestLikeToggler(t *testing.T) {
	target := RandomId()
	caller := RandomId()
	t.Run("target is not already liked - like it", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return false, nil
		}
		t.Run("happy case", func(t *testing.T) {
			like := func(targetId string, liker core_values.UserId) error {
				if targetId == target && liker == caller {
					return nil
				}
				panic("unexpected args")
			}
			err := service.NewLikeToggler(checkLiked, like, nil)(target, caller)
			AssertNoError(t, err)
		})
		t.Run("error case - liking throws", func(t *testing.T) {
			like := func(string, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewLikeToggler(checkLiked, like, nil)(target, caller)
			AssertSomeError(t, err)
		})
	})
	t.Run("target is already liked - unlike it", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return true, nil
		}
		t.Run("happy case", func(t *testing.T) {
			unlike := func(targetId string, unliker core_values.UserId) error {
				if targetId == target && unliker == caller {
					return nil
				}
				panic("unexpected args")
			}
			err := service.NewLikeToggler(checkLiked, nil, unlike)(target, caller)
			AssertNoError(t, err)
		})
		t.Run("error case - unliking throws", func(t *testing.T) {
			unlike := func(string, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewLikeToggler(checkLiked, nil, unlike)(target, caller)
			AssertSomeError(t, err)
		})
	})
	t.Run("checking if target is liked throws", func(t *testing.T) {
		likeChecker := func(targetId string, liker core_values.UserId) (bool, error) {
			if targetId == target && liker == caller {
				return false, RandomError()
			}
			panic("unexpected args")
		}
		err := service.NewLikeToggler(likeChecker, nil, nil)(target, caller)
		AssertSomeError(t, err)
	})
}
