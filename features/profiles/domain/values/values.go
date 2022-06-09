package values

import (
	"core/ref"
	"profiles/domain/entities"
)

type ProfileUpdateData struct {
	About string
}

type AvatarData struct {
	Data ref.Ref[[]byte]
}

type AvatarPath struct {
	Path string
}

type NewProfile struct {
	entities.Profile
}
