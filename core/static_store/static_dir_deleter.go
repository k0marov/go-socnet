package static_store

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"os"
	"path/filepath"
)

// DirDeleter os.RemoveAll implements this
type DirDeleter = func(dir string) error

func NewStaticDirDeleter(deleteDir DirDeleter) StaticDirDeleter {
	return func(dir core_values.StaticPath) error {
		fullDir := filepath.Join(StaticDir, dir)
		err := deleteDir(fullDir)
		if err != nil {
			return fmt.Errorf("while deleting a static dir (%v) : %w", fullDir, err)
		}
		return nil
	}
}

func NewStaticDirDeleterImpl() StaticDirDeleter {
	return NewStaticDirDeleter(os.RemoveAll)
}
