package contexters

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

type ProfileContextAdder func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error)

func NewProfileContextAdder(getContext likeable_contexters.OwnLikeContextGetter) ProfileContextAdder {
	return func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error) {
		context, err := getContext(profile.Id, profile.Id, caller)
		if err != nil {
			return entities.ContextedProfile{}, core_err.Rethrow("getting context for profile", err)
		}
		contextedProfile := entities.ContextedProfile{
			Profile:        profile,
			OwnLikeContext: context,
		}
		return contextedProfile, nil
	}
}
