package entities

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	profile "github.com/k0marov/socnet/features/profiles/domain/entities"
)

type Post struct {
	Id     values.PostId
	Author profile.ContextedProfile
	Text   string
	Images core_values.ImageUrl
}
