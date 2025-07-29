package workdir

import (
	"errors"
	"github.com/otiai10/copy"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type WorkDir struct {
	root string
}

func InitEmptyWorkDir() *WorkDir {
	dir, err := os.MkdirTemp("", "workdir")
	if err != nil {
		panic(err)
	}
	return &WorkDir{root: dir}
}

func (wd *WorkDir) CreateFile(path string) error {
	fullPath := filepath.Join(wd.root, path)
	fatherPath := filepath.Dir(fullPath)

	if err := os.MkdirAll(fatherPath, 0755); err != nil {
		return err
	}

	_, err := os.Create(fullPath)
	return err
}

func (wd *WorkDir) ListFilesRoot() []string {
	entries, err := os.ReadDir(wd.root)
	if err != nil {
		return nil
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}
