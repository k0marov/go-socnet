package file_storage_test

import (
	"github.com/k0marov/go-socnet/core/general/core_values/ref"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/store/file_storage"
)

func TestAvatarFileCreator(t *testing.T) {
	t.Run("should forward the call to static file creator with proper args", func(t *testing.T) {
		tUserId := RandomString()
		tData := []byte(RandomString())
		tDataRef, _ := ref.NewRef(&tData)

		expectedDir := file_storage.ProfilePrefix + tUserId
		expectedName := file_storage.AvatarFileName

		wantPath := RandomString()
		wantErr := RandomError()

		staticFileCreator := func(data ref.Ref[[]byte], dir, name string) (string, error) {
			if data == tDataRef && dir == expectedDir && name == expectedName {
				return wantPath, wantErr
			}
			panic("called with unexpected args")
		}
		sut := file_storage.NewAvatarFileCreator(staticFileCreator)
		path, err := sut(tDataRef, tUserId)
		AssertError(t, err, wantErr)
		Assert(t, path, wantPath, "returned path")
	})
}
