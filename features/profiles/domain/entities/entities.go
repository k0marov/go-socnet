package entities

type Profile struct {
	Id         string
	Username   string
	About      string
	AvatarPath string
}

type DetailedProfile struct {
	Profile
	// ...
}
