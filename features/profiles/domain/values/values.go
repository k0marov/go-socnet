package values

import (
	"core/ref"
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
	Data ref.Ref[[]byte]
}

type AvatarPath struct {
	Path string
}
