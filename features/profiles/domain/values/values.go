package values

import "core/ref"

type ProfileUpdateData struct {
	About string
}

type AvatarData struct {
	Data ref.Ref[[]byte]
}

type AvatarURL struct {
	Url string
}
