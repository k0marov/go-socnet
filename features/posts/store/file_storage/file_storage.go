package file_storage

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/static_store"
	"github.com/k0marov/socnet/features/posts/domain/values"
)

type PostImageFilesCreator = func(values.PostId, core_values.UserId, []core_values.FileData) ([]core_values.StaticFilePath, error)
type PostFilesDeleter = func(values.PostId, core_values.UserId) error

func NewPostImageFilesCreator(createFile static_store.StaticFileCreator) PostImageFilesCreator {
	return func(post values.PostId, author core_values.UserId, images []core_values.FileData) ([]core_values.StaticFilePath, error) {
		panic("unimplemented")
	}
}

func NewPostFilesDeleter(deleteDir static_store.StaticDirDeleter) PostFilesDeleter {
	return func(post values.PostId, author core_values.UserId) error {
		panic("unimplemented")
	}
}
