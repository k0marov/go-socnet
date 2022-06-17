package static_store_test

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/k0marov/socnet/core/ref"
	"github.com/k0marov/socnet/core/static_store"
	. "github.com/k0marov/socnet/core/test_helpers"
)

func TestStaticFileCreator(t *testing.T) {
	tData := []byte(RandomString())
	tDataRef, _ := ref.NewRef(&tData)
	tDir := RandomString()
	tFilename := RandomString()
	wantDir := filepath.Join(static_store.StaticDir, tDir)
	wantFullPath := filepath.Join(wantDir, tFilename)
	wantPath := filepath.Join(tDir, tFilename)

	t.Run("should create directory", func(t *testing.T) {
		recursiveDirCreator := func(path string, perm fs.FileMode) error {
			if path == wantDir && perm == 0777 {
				return nil
			}
			panic("called with unexpected arguments")
		}
		t.Run("happy case, should write to the file", func(t *testing.T) {
			t.Run("happy case", func(t *testing.T) {
				writeFile := func(path string, data []byte, perm fs.FileMode) error {
					if path == wantFullPath && reflect.DeepEqual(data, tData) && perm == 0777 {
						return nil
					}
					panic("called with unexpected arguments")
				}
				sut := static_store.NewStaticFileCreator(recursiveDirCreator, writeFile) // now nil
				gotPath, err := sut(tDataRef, tDir, tFilename)
				AssertNoError(t, err)
				Assert(t, gotPath, wantPath, "returned path")
			})
			t.Run("error case - writing to the file throws", func(t *testing.T) {
				writeFile := func(string, []byte, fs.FileMode) error {
					return RandomError()
				}
				sut := static_store.NewStaticFileCreator(recursiveDirCreator, writeFile)
				_, err := sut(tDataRef, tDir, tFilename)
				AssertSomeError(t, err)
			})
		})
		t.Run("error case - mkdirAll throws", func(t *testing.T) {
			recursiveDirCreator := func(path string, perm fs.FileMode) error {
				return RandomError()
			}
			sut := static_store.NewStaticFileCreator(recursiveDirCreator, nil) // writefile shouldn't be called, so it's nil
			_, err := sut(tDataRef, tDir, tFilename)
			AssertSomeError(t, err)
		})
	})
}
