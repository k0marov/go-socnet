package service_test

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable/service"
	. "github.com/k0marov/socnet/core/test_helpers"
	"testing"
)

func TestLikeToggler(t *testing.T) {
	targetId := RandomId()
	owner := RandomId()
	caller := RandomId()
	t.Run("error case - liking entity that belongs to you", func(t *testing.T) {
		ownerAndLiker := RandomId()
		err := service.NewLikeToggler(nil, nil, nil)(targetId, ownerAndLiker, ownerAndLiker)
		AssertError(t, err, client_errors.LikingYourself)
	})
	t.Run("target is not already liked - like it", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return false, nil
		}
		t.Run("happy case", func(t *testing.T) {
			like := func(target string, liker core_values.UserId) error {
				if target == targetId && liker == caller {
					return nil
				}
				panic("unexpected args")
			}
			err := service.NewLikeToggler(checkLiked, like, nil)(targetId, owner, caller)
			AssertNoError(t, err)
		})
		t.Run("error case - liking throws", func(t *testing.T) {
			like := func(string, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewLikeToggler(checkLiked, like, nil)(targetId, owner, caller)
			AssertSomeError(t, err)
		})
	})
	t.Run("target is already liked - unlike it", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return true, nil
		}
		t.Run("happy case", func(t *testing.T) {
			unlike := func(target string, unliker core_values.UserId) error {
				if target == targetId && unliker == caller {
					return nil
				}
				panic("unexpected args")
			}
			err := service.NewLikeToggler(checkLiked, nil, unlike)(targetId, owner, caller)
			AssertNoError(t, err)
		})
		t.Run("error case - unliking throws", func(t *testing.T) {
			unlike := func(string, core_values.UserId) error {
				return RandomError()
			}
			err := service.NewLikeToggler(checkLiked, nil, unlike)(targetId, owner, caller)
			AssertSomeError(t, err)
		})
	})
	t.Run("checking if target is liked throws", func(t *testing.T) {
		likeChecker := func(target string, liker core_values.UserId) (bool, error) {
			if target == targetId && liker == caller {
				return false, RandomError()
			}
			panic("unexpected args")
		}
		err := service.NewLikeToggler(likeChecker, nil, nil)(targetId, owner, caller)
		AssertSomeError(t, err)
	})
}
