package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/profiles/domain/contexters"
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"testing"
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
