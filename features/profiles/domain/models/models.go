package models

import "github.com/k0marov/go-socnet/core/core_values"

type ProfileModel struct {
	Id         core_values.UserId
	Username   string
	About      string
	AvatarPath string
}
