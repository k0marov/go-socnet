package static_file_creator

import (
	"core/ref"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type StaticFileCreator = func(data ref.Ref[[]byte], dir, filename string) (string, error)

const StaticDir = "static/"

// os.MkdirAll implements this
type RecursiveDirCreator = func(path string, perm fs.FileMode) error

// os.WriteFile implements this
type FileCreator = func(name string, data []byte, perm fs.FileMode) error

func NewStaticFileCreator(mkdirAll RecursiveDirCreator, writeFile FileCreator) StaticFileCreator {
	return func(data ref.Ref[[]byte], dir, filename string) (string, error) {
		fullDir := filepath.Join(StaticDir, dir)
		err := mkdirAll(fullDir, 0777)
		if err != nil {
			return "", fmt.Errorf("error while creating a new directory: %w", err)
		}
		path := filepath.Join(fullDir, filename)
		err = writeFile(path, data.Value(), 0777)
		if err != nil {
			return "", fmt.Errorf("error while writing to a file: %w", err)
		}
		return path, nil
	}
}

func NewStaticFileCreatorImpl() StaticFileCreator {
	return NewStaticFileCreator(os.MkdirAll, os.WriteFile)
}
