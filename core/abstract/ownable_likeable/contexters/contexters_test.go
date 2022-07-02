package contexters_test

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
)

func TestLikeableContextAdder(t *testing.T) {
	target := RandomId()
	caller := RandomId()

	runCase := func(isLiked, isMine bool) {
		t.Run(fmt.Sprintf("isLiked = %v, isMine = %v", isLiked, isMine), func(t *testing.T) {
			checkLiked := func(targetId string, callerId core_values.UserId) (bool, error) {
				if targetId == target && callerId == caller {
					return isLiked, nil
				}
				panic("unexpected args")
			}
			t.Run("error case - getting isLiked throws", func(t *testing.T) {
				checkLiked := func(string, core_values.UserId) (bool, error) {
					return false, RandomError()
				}
				_, err := contexters.NewOwnLikeContextGetter(checkLiked, nil)(target, caller)
				AssertSomeError(t, err)
			})
			t.Run("happy case", func(t *testing.T) {
				var owner core_values.UserId
				if isMine {
					owner = caller
				} else {
					owner = RandomId()
				}
				getOwner := func(targetId string) (core_values.UserId, error) {
					if targetId == target {
						return owner, nil
					}
					panic("unexpected args")
				}
				gotContext, err := contexters.NewOwnLikeContextGetter(checkLiked, getOwner)(target, caller)
				AssertNoError(t, err)
				wantContext := contexters.OwnLikeContext{
					IsLiked: isLiked,
					IsMine:  isMine,
				}
				Assert(t, gotContext, wantContext, "returned context")
			})

			t.Run("error case - getting owner throws", func(t *testing.T) {
				getOwner := func(string) (core_values.UserId, error) {
					return "", RandomError()
				}
				_, err := contexters.NewOwnLikeContextGetter(checkLiked, getOwner)(target, caller)
				AssertSomeError(t, err)
			})
		})
	}

	runCase(true, true)
	runCase(true, false)
	runCase(false, true)
	runCase(false, false)
}
