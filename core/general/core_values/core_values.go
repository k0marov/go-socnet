package core_values

import (
	"github.com/k0marov/go-socnet/core/general/core_values/ref"
)

type UserId = string

type FileData = ref.Ref[[]byte]

type FileURL = string
type StaticPath = string
