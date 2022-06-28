package static_store

import (
	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/ref"
	"log"
	"os"
)

type (
	StaticFileCreator = func(data ref.Ref[[]byte], dir, filename string) (core_values.StaticPath, error)
	StaticDirDeleter  = func(dir core_values.StaticPath) error
)

const StaticDir = "static/"

var StaticHost = getStaticHostStr()

func PathToURL(path core_values.StaticPath) core_values.FileURL {
	if path == "" {
		return ""
	}
	return StaticHost + "/" + path
}

func getStaticHostStr() string {
	const staticHostEnv = "SOCIO_STATIC_HOST"
	host, exists := os.LookupEnv(staticHostEnv)
	if !exists {
		log.Fatalf(`Environment variable %s is not set.
If this is a test, just set the environment variable to a dummy string.
If this is in production, set this environment variable to point to the URL from which the static directory can be accessed.`, staticHostEnv)
	}
	return host
}
