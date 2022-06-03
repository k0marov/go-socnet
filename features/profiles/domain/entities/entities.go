package entities

type Profile struct {
	Id       string
	Username string
	About    string
}

type DetailedProfile struct {
	Profile
	// ...
}
