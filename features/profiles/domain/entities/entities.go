package entities

import "github.com/k0marov/socnet/core/core_values"

type Profile struct {
	Id         core_values.UserId
	Username   string
	About      string
	AvatarPath string
	Follows    int
	Followers  int
	// IsFollowedByCaller bool
}

type ContextedProfile struct {
	Profile
	IsFollowedByCaller bool
}
