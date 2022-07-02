package core_entities

import auth "github.com/k0marov/golang-auth"

type User struct {
	Id       string
	Username string
}

func UserFromAuth(authUser auth.User) User {
	return User{
		Id:       authUser.Id,
		Username: authUser.Username,
	}
}
