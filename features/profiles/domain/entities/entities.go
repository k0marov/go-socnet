package entities

import "core/core_values"

type Profile struct {
	Id         core_values.UserId
	Username   string
	About      string
	AvatarPath string
	Follows    int
	Followers  int
}

type DetailedProfile struct {
	Profile
	FollowsProfiles []Profile
}
