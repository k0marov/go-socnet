package values

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type PostId = string

type NewPostData struct {
	Author core_values.UserId
	Text   string
	Images []PostImageFile
}

type PostImageFile struct {
	File  core_values.FileData
	Index int
}

type PostImage struct {
	URL   core_values.FileURL
	Index int
}
