package image_decoder_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/k0marov/socnet/core/image_decoder"
	. "github.com/k0marov/socnet/core/test_helpers"
)

func TestImageDecoderImpl(t *testing.T) {
	readFixture := func(fileName string) []byte {
		file, err := os.Open(filepath.Join("testdata", fileName))
		if err != nil {
			t.Fatalf(fmt.Sprintf("error while opening fixture file: %v", err))
		}
		defer file.Close()
		contents, _ := io.ReadAll(file)
		return contents
	}

	cases := []struct {
		fixtureName string
		wantImg     image_decoder.Image
		shouldErr   bool
	}{
		{"test_avatar.jpg", image_decoder.Image{Width: 640, Height: 640}, false},
		{"test_non_square_avatar.png", image_decoder.Image{Width: 640, Height: 480}, false},
		{"test_js_injection.js", image_decoder.Image{}, true},
		{"test_text.txt", image_decoder.Image{}, true},
	}
	for _, c := range cases {
		t.Run(c.fixtureName, func(t *testing.T) {
			fileData := readFixture(c.fixtureName)
			img, err := image_decoder.ImageDecoderImpl(fileData)
			Assert(t, img, c.wantImg, "returned image")
			if c.shouldErr {
				AssertSomeError(t, err)
			} else {
				AssertNoError(t, err)
			}
		})
	}
}
