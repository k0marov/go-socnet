package validators_test

import (
	"core/client_errors"
	"core/image_decoder"
	. "core/test_helpers"
	"fmt"
	"profiles/domain/validators"
	"profiles/domain/values"
	"strings"
	"testing"
)

func TestProfileUpdateValidator(t *testing.T) {
	cases := []struct {
		profileUpdate values.ProfileUpdateData
		ok            bool
		err           client_errors.ClientError
	}{
		{values.ProfileUpdateData{About: "abcdfeg"}, true, client_errors.ClientError{}},
		{values.ProfileUpdateData{About: ""}, true, client_errors.ClientError{}},
		{values.ProfileUpdateData{About: strings.Repeat("abc", 100)}, false, client_errors.AboutTooLong},
	}
	sut := validators.NewProfileUpdateValidator()
	for _, c := range cases {
		t.Run(c.profileUpdate.About, func(t *testing.T) {
			gotErr, gotOk := sut(c.profileUpdate)
			Assert(t, gotOk, c.ok, "validation result")
			AssertError(t, gotErr, c.err)
		})
	}
}

func TestAvatarValidator(t *testing.T) {
	goodAvatar := []byte(RandomString())
	nonSquareAvatar := []byte(RandomString())
	jsInjectionAvatar := []byte(RandomString())

	imageDecoder := func(fileContents *[]byte) (image_decoder.Image, error) {
		if fileContents == &goodAvatar {
			return image_decoder.Image{Width: 10, Height: 10}, nil
		} else if fileContents == &nonSquareAvatar {
			return image_decoder.Image{Width: 10, Height: 20}, nil
		} else if fileContents == &jsInjectionAvatar {
			return image_decoder.Image{}, RandomError()
		}
		panic(fmt.Sprintf("called with incorrect arguments, fileContents=%v", fileContents))
	}
	sut := validators.NewAvatarValidator(imageDecoder)

	cases := []struct {
		avatar values.AvatarData
		ok     bool
		err    client_errors.ClientError
	}{
		{values.AvatarData{Data: &goodAvatar, FileName: RandomString()}, true, client_errors.ClientError{}},
		{values.AvatarData{Data: &jsInjectionAvatar, FileName: RandomString()}, false, client_errors.NonImageAvatar},
		{values.AvatarData{Data: &nonSquareAvatar, FileName: RandomString()}, false, client_errors.NonSquareAvatar},
	}

	for _, c := range cases {
		gotErr, gotOk := sut(c.avatar)
		Assert(t, gotOk, c.ok, "result of validation")
		AssertError(t, gotErr, c.err)
	}
}
