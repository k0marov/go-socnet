package file_storage_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/store/file_storage"
	profiles "github.com/k0marov/socnet/features/profiles/store/file_storage"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPostImageFilesCreator(t *testing.T) {
	post := RandomString()
	author := RandomString()
	images := []core_values.FileData{RandomFileData(), RandomFileData(), RandomFileData()}
	paths := []core_values.StaticFilePath{RandomString(), RandomString(), RandomString()}
	t.Run("happy case", func(t *testing.T) {
		var calledWithImages []core_values.FileData
		wantDir := filepath.Join(profiles.ProfilePrefix+author, file_storage.PostPrefix+post)
		createFile := func(file core_values.FileData, dir string, filename string) (core_values.StaticFilePath, error) {
			if filename == file_storage.ImagePrefix+strconv.Itoa(len(calledWithImages)+1) &&
				dir == wantDir {
				calledWithImages = append(calledWithImages, file)
				return paths[len(calledWithImages)-1], nil
			}
			panic("unexpected args")
		}
		gotPaths, err := file_storage.NewPostImageFilesCreator(createFile)(post, author, images)
		AssertNoError(t, err)
		Assert(t, gotPaths, paths, "returned paths")
		Assert(t, calledWithImages, images, "fileData args to createFile")
	})
	t.Run("error case", func(t *testing.T) {
		createFile := func(core_values.FileData, string, string) (core_values.StaticFilePath, error) {
			return "", RandomError()
		}
		_, err := file_storage.NewPostImageFilesCreator(createFile)(post, author, images)
		AssertSomeError(t, err)
	})
}

func TestPostFilesDeleter(t *testing.T) {

}
