package entities

import (
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	"github.com/k0marov/socnet/features/profiles/domain/models"
)

type Profile struct {
	models.ProfileModel
	Follows   int
	Followers int
}

type ContextedProfile struct {
	Profile
	likeable_contexters.LikeableContext
}
