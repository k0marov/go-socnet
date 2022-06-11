package values

import (
	"core/ref"
)

type UserId = string

type ProfileUpdateData struct {
	About string
}

type AvatarData struct {
	Data ref.Ref[[]byte]
}

type AvatarPath struct {
	Path string
}
