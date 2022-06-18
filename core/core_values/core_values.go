package core_values

import "github.com/k0marov/socnet/core/ref"

type UserId = string

type FileData = ref.Ref[[]byte]

type FileURL = string
type StaticPath = string
