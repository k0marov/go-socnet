package validators_test

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/image_decoder"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/domain/validators"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"reflect"
	"strings"
	"testing"
)

func TestPostValidator(t *testing.T) {
	t.Run("Text validation", func(t *testing.T) {
		decoder := func([]byte) (image_decoder.Image, error) {
			return image_decoder.Image{Height: 123, Width: 345}, nil
		}
		cases := []struct {
			text        string
			expectedErr error
		}{
			{"", nil},
			{"some short text", nil},
			{strings.Repeat("looong", 300), client_errors.TextTooLong},
		}
		for _, testCase := range cases {
			t.Run(testCase.text, func(t *testing.T) {
				newPost := values.NewPostData{
					Author: RandomString(),
					Text:   testCase.text,
					Images: nil,
				}
				gotErr, ok := validators.NewPostValidator(decoder)(newPost)
				if testCase.expectedErr == nil {
					AssertError(t, gotErr, client_errors.ClientError{})
					Assert(t, ok, true, "returned 'ok' value")
				} else {
					AssertError(t, gotErr, testCase.expectedErr)
					Assert(t, ok, false, "returned 'ok' value")
				}
			})
		}
	})
	t.Run("Images validation", func(t *testing.T) {
		newPost := values.NewPostData{
			Author: RandomString(),
			Text:   "",
			Images: []values.PostImageFile{{RandomFileData(), 1}, {RandomFileData(), 2}},
		}
		isCorrectImage := func(image []byte) bool {
			for _, postImage := range newPost.Images {
				if reflect.DeepEqual(postImage.File.Value(), image) {
					return true
				}
			}
			return false
		}
		t.Run("happy case", func(t *testing.T) {
			imagesChecked := 0
			decoder := func(image []byte) (image_decoder.Image, error) {
				if isCorrectImage(image) {
					imagesChecked++
					return image_decoder.Image{Height: 420, Width: 840}, nil
				}
				panic("unexpected args")
			}
			clientErr, ok := validators.NewPostValidator(decoder)(newPost)
			Assert(t, ok, true, "ok is true")
			AssertError(t, clientErr, client_errors.ClientError{})
			Assert(t, imagesChecked, len(newPost.Images), "amount of checked images")
		})
		t.Run("error case", func(t *testing.T) {
			decoder := func([]byte) (image_decoder.Image, error) {
				return image_decoder.Image{}, RandomError()
			}
			clientErr, ok := validators.NewPostValidator(decoder)(newPost)
			Assert(t, ok, false, "ok is false")
			AssertError(t, clientErr, client_errors.InvalidImage)
		})
	})
}
