package responses

import "github.com/k0marov/socnet/features/profiles/domain/entities"

type ProfileResponse struct {
	Id         string
	Username   string
	About      string
	AvatarPath string
	Follows    int
	Followers  int
	IsMine     bool
	IsFollowed bool
}

func NewProfileResponse(profile entities.ContextedProfile) ProfileResponse {
	return ProfileResponse{
		Id:         profile.Id,
		Username:   profile.Username,
		About:      profile.About,
		AvatarPath: profile.AvatarPath,
		Follows:    profile.Follows,
		Followers:  profile.Followers,
		IsMine:     profile.IsMine,
		IsFollowed: profile.IsLiked,
	}
}
