package file_storage_test

import (
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/k0marov/go-socnet/core/core_values"
	. "github.com/k0marov/go-socnet/core/test_helpers"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	"github.com/k0marov/go-socnet/features/posts/store/file_storage"
	profiles "github.com/k0marov/go-socnet/features/profiles/store/file_storage"
)

func TestPostImageFilesCreator(t *testing.T) {
	post := RandomString()
	author := RandomString()
	images := []values.PostImageFile{{RandomFileData(), 1}, {RandomFileData(), 2}}
	paths := []core_values.StaticPath{RandomString(), RandomString()}
	t.Run("happy case", func(t *testing.T) {
		wantDir := filepath.Join(profiles.ProfilePrefix+author, file_storage.PostPrefix+post)
		filesStored := 0
		createFile := func(file core_values.FileData, dir string, filename string) (core_values.StaticPath, error) {
			wantFilename := file_storage.ImagePrefix + strconv.Itoa(filesStored+1)
			if filename == wantFilename && dir == wantDir && reflect.DeepEqual(file, images[filesStored].File) {
				filesStored++
				return paths[filesStored-1], nil
			}
			panic("unexpected args")
		}
		gotPaths, err := file_storage.NewPostImageFilesCreator(createFile)(post, author, images)
		AssertNoError(t, err)
		Assert(t, gotPaths, paths, "returned paths")
		Assert(t, filesStored, len(images), "number of stored files")
	})
	t.Run("error case", func(t *testing.T) {
		createFile := func(core_values.FileData, string, string) (core_values.StaticPath, error) {
			return "", RandomError()
		}
		_, err := file_storage.NewPostImageFilesCreator(createFile)(post, author, images)
		AssertSomeError(t, err)
	})
}

func TestPostFilesDeleter(t *testing.T) {
	post := RandomString()
	author := RandomString()
	t.Run("happy case", func(t *testing.T) {
		wantDir := filepath.Join(profiles.GetProfileDir(author), file_storage.PostPrefix+post)
		deleteDir := func(dir string) error {
			if dir == wantDir {
				return nil
			}
			panic("unexpected args")
		}
		err := file_storage.NewPostFilesDeleter(deleteDir)(post, author)
		AssertNoError(t, err)
	})
	t.Run("error case", func(t *testing.T) {
		deleteDir := func(string) error {
			return RandomError()
		}
		err := file_storage.NewPostFilesDeleter(deleteDir)(post, author)
		AssertSomeError(t, err)
	})
}
