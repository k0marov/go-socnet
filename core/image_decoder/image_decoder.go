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

type ImageDecoder = func(fileData *[]byte) (Image, error)

func ImageDecoderImpl(fileData *[]byte) (Image, error) {
	imageConfig, _, err := image.DecodeConfig(bytes.NewReader(*fileData)) // TODO: maybe remove dereferencing for performance?
	if err != nil {
		return Image{}, fmt.Errorf("error while decoding image: %w", err)
	}
	return Image{Width: imageConfig.Width, Height: imageConfig.Height}, nil
}
