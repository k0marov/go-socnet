package service_test

import (
	"github.com/k0marov/go-socnet/core/abstract/deletable/service"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestDeleter(t *testing.T) {
	target := RandomId()
	owner := RandomId()

	caller := owner

	getOwner := func(targetId string) (core_values.UserId, error) {
		if targetId == target {
			return owner, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - caller is not owner", func(t *testing.T) {
		err := service.NewDeleter(getOwner, nil)(target, RandomId())
		AssertError(t, err, client_errors.InsufficientPermissions)
	})
	t.Run("error case - getting owner throws", func(t *testing.T) {
		getOwner := func(string) (core_values.UserId, error) {
			return owner, RandomError()
		}
		err := service.NewDeleter(getOwner, nil)(target, caller)
		AssertSomeError(t, err)
	})

	delete := func(targetId string) error {
		if targetId == target {
			return nil
		}
		panic("unexpected args")
	}
	t.Run("error case - deleting throws", func(t *testing.T) {
		delete := func(string) error {
			return RandomError()
		}
		err := service.NewDeleter(getOwner, delete)(target, caller)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		err := service.NewDeleter(getOwner, delete)(target, caller)
		AssertNoError(t, err)
	})
}
