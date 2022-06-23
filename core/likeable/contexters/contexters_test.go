package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable/contexters"
	. "github.com/k0marov/socnet/core/test_helpers"
	"testing"
)

type TestTarget struct {
	id    string
	owner core_values.UserId
}

func (t TestTarget) GetId() string {
	return t.id
}
func (t TestTarget) GetOwner() core_values.UserId {
	return t.owner
}

func TestLikeableContextAdder(t *testing.T) {
	target := TestTarget{RandomId(), RandomId()}
	wantContext := contexters.LikeableContext{
		IsLiked: RandomBool(),
		IsMine:  RandomBool(),
	}
	var caller core_values.UserId
	if wantContext.IsMine {
		caller = target.owner
	} else {
		caller = RandomId()
	}
	t.Run("happy case", func(t *testing.T) {
		checkLiked := func(targetId string, callerId core_values.UserId) (bool, error) {
			if targetId == target.id && callerId == caller {
				return wantContext.IsLiked, nil
			}
			panic("unexpected args")
		}
		context, err := contexters.NewLikeableContextGetter(checkLiked)(target, caller)
		AssertNoError(t, err)
		Assert(t, context, wantContext, "the returned context")
	})
	t.Run("error case - checking if target is liked throws", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return false, RandomError()
		}
		_, err := contexters.NewLikeableContextGetter(checkLiked)(target, caller)
		AssertSomeError(t, err)
	})

}
