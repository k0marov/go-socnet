package static_store

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/ref"
)

type (
	StaticFileCreator = func(data ref.Ref[[]byte], dir, filename string) (core_values.StaticFilePath, error)
	StaticDirDeleter  = func(dir string) error
)

const StaticDir = "static/"
const StaticHost = "static.example.com"

func PathToURL(path core_values.StaticFilePath) core_values.FileURL {
	return StaticHost + "/" + path
}
