package static_store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/k0marov/socnet/core/ref"
)

// RecursiveDirCreator os.MkdirAll implements this
type RecursiveDirCreator = func(path string, perm fs.FileMode) error

// FileCreator os.WriteFile implements this
type FileCreator = func(name string, data []byte, perm fs.FileMode) error

func NewStaticFileCreator(mkdirAll RecursiveDirCreator, writeFile FileCreator) StaticFileCreator {
	return func(data ref.Ref[[]byte], dir, filename string) (string, error) {
		fullDir := filepath.Join(StaticDir, dir)
		err := mkdirAll(fullDir, 0777)
		if err != nil {
			return "", fmt.Errorf("error while creating a new directory: %w", err)
		}
		fullPath := filepath.Join(fullDir, filename)
		err = writeFile(fullPath, data.Value(), 0777)
		if err != nil {
			return "", fmt.Errorf("error while writing to a file: %w", err)
		}
		path := filepath.Join(dir, filename)
		return path, nil
	}
}

func NewStaticFileCreatorImpl() StaticFileCreator {
	return NewStaticFileCreator(os.MkdirAll, os.WriteFile)
}
