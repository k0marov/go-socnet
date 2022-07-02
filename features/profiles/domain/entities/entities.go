package entities

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_values"
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
	likeable_contexters.OwnLikeContext
}
