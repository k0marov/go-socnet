package validators_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/domain/validators"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/image_decoder"
	"github.com/k0marov/go-socnet/core/ref"
	. "github.com/k0marov/go-socnet/core/test_helpers"
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

func makeRefWithoutCheck(data *[]byte) ref.Ref[[]byte] {
	ref, err := ref.NewRef(data)
	if err != nil {
		panic("ref data was nil")
	}
	return ref
}

func TestAvatarValidator(t *testing.T) {
	goodAvatar := []byte(RandomString())
	nonSquareAvatar := []byte(RandomString())
	jsInjectionAvatar := []byte(RandomString())

	imageDecoder := func(fileContents []byte) (image_decoder.Image, error) {
		if reflect.DeepEqual(fileContents, goodAvatar) {
			return image_decoder.Image{Width: 10, Height: 10}, nil
		} else if reflect.DeepEqual(fileContents, nonSquareAvatar) {
			return image_decoder.Image{Width: 10, Height: 20}, nil
		} else if reflect.DeepEqual(fileContents, jsInjectionAvatar) {
			return image_decoder.Image{}, RandomError()
		}
		panic(fmt.Sprintf("called with unexpected arguments, fileContents=%v", fileContents))
	}
	sut := validators.NewAvatarValidator(imageDecoder)

	cases := []struct {
		avatar values.AvatarData
		ok     bool
		err    client_errors.ClientError
	}{
		{values.AvatarData{Data: makeRefWithoutCheck(&goodAvatar)}, true, client_errors.ClientError{}},
		{values.AvatarData{Data: makeRefWithoutCheck(&nonSquareAvatar)}, false, client_errors.NonSquareAvatar},
		{values.AvatarData{Data: makeRefWithoutCheck(&jsInjectionAvatar)}, false, client_errors.InvalidImage},
	}

	for _, c := range cases {
		gotErr, gotOk := sut(c.avatar)
		Assert(t, gotOk, c.ok, "result of validation")
		AssertError(t, gotErr, c.err)
	}
}
