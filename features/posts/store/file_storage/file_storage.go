package file_storage

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
)

type PostImageFilesCreator = func(values.PostId, core_values.UserId, []core_values.FileData) ([]core_values.StaticFilePath, error)
