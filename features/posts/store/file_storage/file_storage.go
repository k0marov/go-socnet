package file_storage

import (
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/static_store"
	"path/filepath"
	"strconv"

	"github.com/k0marov/go-socnet/features/posts/domain/values"

	profiles "github.com/k0marov/go-socnet/features/profiles/store/file_storage"
)

const PostPrefix = "post_"
const ImagePrefix = "image_"

type PostImageFilesCreator = func(values.PostId, core_values.UserId, []values.PostImageFile) ([]core_values.StaticPath, error)
type PostFilesDeleter = func(values.PostId, core_values.UserId) error

func NewPostImageFilesCreator(createFile static_store.StaticFileCreator) PostImageFilesCreator {
	return func(post values.PostId, author core_values.UserId, images []values.PostImageFile) (paths []core_values.StaticPath, err error) {
		dir := filepath.Join(profiles.ProfilePrefix+author, PostPrefix+post)
		for _, image := range images {
			filename := ImagePrefix + strconv.Itoa(image.Index)
			path, err := createFile(image.File, dir, filename)
			if err != nil {
				return paths, core_err.Rethrow("storing a file", err)
			}
			paths = append(paths, path)
		}
		return
	}
}

func NewPostFilesDeleter(deleteDir static_store.StaticDirDeleter) PostFilesDeleter {
	return func(post values.PostId, author core_values.UserId) error {
		err := deleteDir(GetPostDir(post, author))
		if err != nil {
			return core_err.Rethrow("deleting the post directory", err)
		}
		return nil
	}
}

func GetPostDir(post values.PostId, author core_values.UserId) string {
	return filepath.Join(profiles.GetProfileDir(author), PostPrefix+post)
}
