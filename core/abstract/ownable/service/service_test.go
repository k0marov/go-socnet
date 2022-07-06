package service_test

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable/service"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestOwnerGetter(t *testing.T) {
	target := RandomString()
	owner := RandomString()

	getOwner := func(targetId string) (string, error) {
		if targetId == target {
			return owner, nil
		}
		panic("unexpected args")
	}

	t.Run("error case - getting owner returns ErrNotFound", func(t *testing.T) {
		getOwner := func(string) (string, error) {
			return "", core_err.ErrNotFound
		}
		_, err := service.NewOwnerGetter(getOwner)(target)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - getting owner returns some other error", func(t *testing.T) {
		getOwner := func(string) (string, error) {
			return "", RandomError()
		}
		_, err := service.NewOwnerGetter(getOwner)(target)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		gotOwner, err := service.NewOwnerGetter(getOwner)(target)
		AssertNoError(t, err)
		Assert(t, gotOwner, owner, "returned owner")
	})

}
