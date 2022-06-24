package responses

import "github.com/k0marov/socnet/features/profiles/domain/entities"

type ProfileResponse struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	About      string `json:"about"`
	AvatarPath string `json:"avatar_path"`
	Follows    int    `json:"follows"`
	Followers  int    `json:"followers"`
	IsMine     bool   `json:"is_mine"`
	IsFollowed bool   `json:"is_followed"`
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
