package file_storage

import (
	"io/fs"
	"profiles/data/store"
)

type StaticFileCreator = func(data []byte, dir, filename string)

func NewAvatarFileCreator(fs fs.FS) store.AvatarFileCreator {
	return func(data []byte, belongsToUser string) (string, error) {
		panic("unimplemented")
	}
}
