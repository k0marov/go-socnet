package static_store

import "github.com/k0marov/socnet/core/ref"

type (
	StaticFileCreator = func(data ref.Ref[[]byte], dir, filename string) (string, error)
	StaticDirDeleter  = func(dir string) error
)

const StaticDir = "static/"
