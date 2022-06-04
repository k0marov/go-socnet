package image_decoder

import (
	"bytes"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Image struct {
	Width  int
	Height int
}

type ImageDecoder interface {
	Decode(*[]byte) (Image, error)
}

type ImageDecoderImpl struct{}

func (i ImageDecoderImpl) Decode(data *[]byte) (Image, error) {
	imageConfig, _, err := image.DecodeConfig(bytes.NewReader(*data))
	if err != nil {
		return Image{}, fmt.Errorf("error while decoding image: %w", err)
	}
	return Image{Width: imageConfig.Width, Height: imageConfig.Height}, nil
}
