package models

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type ProfileModel struct {
	Id         core_values.UserId `db:"id"`
	Username   string             `db:"username"`
	About      string             `db:"about"`
	AvatarPath string             `db:"avatarPath"`
}
