package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
)

type Profile struct {
	Id         core_values.UserId
	Username   string
	About      string
	AvatarPath string
	Follows    int
	Followers  int
}

type ContextedProfile struct {
	Profile
	likeable_contexters.LikeableContext
}
