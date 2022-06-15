package values

import (
	"core/core_values"
)

type PostId = string

type NewPostData struct {
	Author core_values.UserId
	Text   string
	Images []core_values.Image
}
