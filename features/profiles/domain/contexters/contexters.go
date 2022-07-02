package contexters

import (
	"fmt"
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_values"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

type ProfileContextAdder func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error)

func NewProfileContextAdder(getContext likeable_contexters.LikeableContextGetter) ProfileContextAdder {
	return func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error) {
		context, err := getContext(profile.Id, profile.Id, caller)
		if err != nil {
			return entities.ContextedProfile{}, fmt.Errorf("while getting context for profile: %w", err)
		}
		contextedProfile := entities.ContextedProfile{
			Profile:         profile,
			LikeableContext: context,
		}
		return contextedProfile, nil
	}
}
