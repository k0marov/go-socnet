package values

import (
	"github.com/k0marov/socnet/core/core_values"
)

type PostId = string

type NewPostData struct {
	Author core_values.UserId
	Text   string
	Images []core_values.FileData
}

type PostImage struct {
	Path  core_values.FileURL
	Index int
}
