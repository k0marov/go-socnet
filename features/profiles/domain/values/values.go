package values

type ProfileUpdateData struct {
	About string
}

type AvatarData struct {
	Data     *[]byte
	FileName string
}

type AvatarURL struct {
	Url string
}
