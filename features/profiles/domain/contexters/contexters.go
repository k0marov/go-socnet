package contexters

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/profiles/domain/entities"
)

type ProfileContextAdder func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error)

func NewProfileContextAdder() ProfileContextAdder {
	return func(profile entities.Profile, caller core_values.UserId) (entities.ContextedProfile, error) {
		panic("unimplemented")
	}
}
