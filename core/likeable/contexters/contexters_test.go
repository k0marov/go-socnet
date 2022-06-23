package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable/contexters"
	. "github.com/k0marov/socnet/core/test_helpers"
	"testing"
)

func TestLikeableContextAdder(t *testing.T) {
	target := RandomId()
	owner := RandomId()
	wantContext := contexters.LikeableContext{
		IsLiked: RandomBool(),
		IsMine:  RandomBool(),
	}
	var caller core_values.UserId
	if wantContext.IsMine {
		caller = owner
	} else {
		caller = RandomId()
	}
	t.Run("happy case", func(t *testing.T) {
		checkLiked := func(targetId string, callerId core_values.UserId) (bool, error) {
			if targetId == target && callerId == caller {
				return wantContext.IsLiked, nil
			}
			panic("unexpected args")
		}
		context, err := contexters.NewLikeableContextGetter(checkLiked)(target, owner, caller)
		AssertNoError(t, err)
		Assert(t, context, wantContext, "the returned context")
	})
	t.Run("error case - checking if target is liked throws", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return false, RandomError()
		}
		_, err := contexters.NewLikeableContextGetter(checkLiked)(target, owner, caller)
		AssertSomeError(t, err)
	})

}
