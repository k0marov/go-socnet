package contexters_test

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/domain/contexters"
	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

func TestProfileContextAdder(t *testing.T) {
	profile := RandomProfile()
	caller := RandomId()
	t.Run("happy case", func(t *testing.T) {
		context := RandomLikeableContext()
		getContext := func(targetId string, owner, callerId core_values.UserId) (likeable_contexters.LikeableContext, error) {
			if targetId == profile.Id && owner == profile.Id && callerId == caller {
				return context, nil
			}
			panic("unexpected args")
		}
		contextedProfile, err := contexters.NewProfileContextAdder(getContext)(profile, caller)
		AssertNoError(t, err)
		wantProfile := entities.ContextedProfile{
			Profile:         profile,
			LikeableContext: context,
		}
		Assert(t, contextedProfile, wantProfile, "returned profile")
	})
	t.Run("error case - getting context throws", func(t *testing.T) {
		getContext := func(targetId string, owner, callerId core_values.UserId) (likeable_contexters.LikeableContext, error) {
			return likeable_contexters.LikeableContext{}, RandomError()
		}
		_, err := contexters.NewProfileContextAdder(getContext)(profile, caller)
		AssertSomeError(t, err)
	})

}
