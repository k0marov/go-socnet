package entities

import "profiles/domain/values"

type Profile struct {
	Id         values.UserId
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
