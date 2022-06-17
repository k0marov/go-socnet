package file_storage

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/static_store"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"path/filepath"
	"strconv"

	profiles "github.com/k0marov/socnet/features/profiles/store/file_storage"
)

const PostPrefix = "post_"
const ImagePrefix = "image_"

type PostImageFilesCreator = func(values.PostId, core_values.UserId, []values.PostImageFile) ([]core_values.StaticFilePath, error)
type PostFilesDeleter = func(values.PostId, core_values.UserId) error

func NewPostImageFilesCreator(createFile static_store.StaticFileCreator) PostImageFilesCreator {
	return func(post values.PostId, author core_values.UserId, images []values.PostImageFile) (paths []core_values.StaticFilePath, err error) {
		dir := filepath.Join(profiles.ProfilePrefix+author, PostPrefix+post)
		for _, image := range images {
			filename := ImagePrefix + strconv.Itoa(image.Index)
			path, err := createFile(image.File, dir, filename)
			if err != nil {
				return paths, fmt.Errorf("while storing a file: %w", err)
			}
			paths = append(paths, path)
		}
		return
	}
}

func NewPostFilesDeleter(deleteDir static_store.StaticDirDeleter) PostFilesDeleter {
	return func(post values.PostId, author core_values.UserId) error {
		dir := filepath.Join(profiles.ProfilePrefix+author, PostPrefix+post)
		err := deleteDir(dir)
		if err != nil {
			return fmt.Errorf("while deleting the post directory: %w", err)
		}
		return nil
	}
}
