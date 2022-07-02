package service_test

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable/service"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestSafeLikeToggler(t *testing.T) {
	target := RandomId()
	owner := RandomId()
	caller := RandomId()

	getOwner := func(targetId string) (core_values.UserId, error) {
		if targetId == target {
			return owner, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - caller is owner", func(t *testing.T) {
		err := service.NewSafeLikeToggler(getOwner, nil)(target, owner)
		AssertError(t, err, client_errors.LikingYourself)
	})
	t.Run("error case - getting author throws", func(t *testing.T) {
		getOwner := func(targetId string) (core_values.UserId, error) {
			return "", RandomError()
		}
		err := service.NewSafeLikeToggler(getOwner, nil)(target, caller)
		AssertSomeError(t, err)
	})
	toggleLike := func(targetId string, callerId core_values.UserId) error {
		if targetId == target && callerId == caller {
			return nil
		}
		panic("unexpected args")
	}
	t.Run("error case - toggling like throws", func(t *testing.T) {
		toggleLike := func(string, core_values.UserId) error {
			return RandomError()
		}
		err := service.NewSafeLikeToggler(getOwner, toggleLike)(target, caller)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		err := service.NewSafeLikeToggler(getOwner, toggleLike)(target, caller)
		AssertNoError(t, err)
	})
}
