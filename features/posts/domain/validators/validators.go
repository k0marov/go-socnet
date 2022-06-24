package validators

import (
	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/image_decoder"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type PostValidator func(newPost values.NewPostData) (client_errors.ClientError, bool)

const MaxTextLength = 1000

func NewPostValidator(decodeImg image_decoder.ImageDecoder) PostValidator {
	return func(newPost values.NewPostData) (client_errors.ClientError, bool) {
		if len(newPost.Text) > MaxTextLength {
			return client_errors.TextTooLong, false
		}
		for _, image := range newPost.Images {
			_, err := decodeImg(image.File.Value())
			if err != nil {
				return client_errors.InvalidImage, false
			}
		}
		return client_errors.ClientError{}, true
	}
}
