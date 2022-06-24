package entities

import (
	"github.com/k0marov/go-socnet/core/core_values"
	likeable_contexters "github.com/k0marov/go-socnet/core/likeable/contexters"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"
)

type Profile struct {
	models.ProfileModel
	AvatarURL core_values.FileURL
	Follows   int
	Followers int
}

type ContextedProfile struct {
	Profile
	likeable_contexters.LikeableContext
}
