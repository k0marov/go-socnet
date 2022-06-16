package values

import (
	"github.com/k0marov/socnet/core/core_values"
)

type ProfileUpdateData struct {
	About string
}
type NewProfile struct {
	Id         string
	Username   string
	About      string
	AvatarPath string
}

type AvatarData struct {
	Data core_values.FileData
}
