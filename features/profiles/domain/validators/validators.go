package validators

import (
	"github.com/k0marov/go-socnet/features/profiles/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/image_decoder"
)

type ProfileUpdateValidator func(values.ProfileUpdateData) (client_errors.ClientError, bool)
type AvatarValidator func(values.AvatarData) (client_errors.ClientError, bool)

const MaxAboutLength = 255

func NewProfileUpdateValidator() ProfileUpdateValidator {
	return func(profileUpdate values.ProfileUpdateData) (client_errors.ClientError, bool) {
		if len(profileUpdate.About) > MaxAboutLength {
			return client_errors.AboutTooLong, false
		}
		return client_errors.ClientError{}, true
	}
}

func NewAvatarValidator(imageDecoder image_decoder.ImageDecoder) AvatarValidator {
	return func(avatar values.AvatarData) (client_errors.ClientError, bool) {
		imageDimensions, err := imageDecoder(avatar.Data.Value())
		if err != nil {
			return client_errors.InvalidImage, false
		}
		if imageDimensions.Height != imageDimensions.Width {
			return client_errors.NonSquareAvatar, false
		}
		return client_errors.ClientError{}, true
	}
}
