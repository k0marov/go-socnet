package static_store

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/ref"
)

type (
	StaticFileCreator = func(data ref.Ref[[]byte], dir, filename string) (core_values.StaticPath, error)
	StaticDirDeleter  = func(dir core_values.StaticPath) error
)

const StaticDir = "static/"
const StaticHost = "static.example.com"

func PathToURL(path core_values.StaticPath) core_values.FileURL {
	return StaticHost + "/" + path
}
