package workdir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

type WorkDir struct {
	root string
}

func NewWorkDir(root string) *WorkDir {
	return &WorkDir{root: root}
}

func InitEmptyWorkDir() *WorkDir {
	dir, err := os.MkdirTemp("", "workdir")
	if err != nil {
		panic(err)
	}
	return &WorkDir{root: dir}
}

func (wd *WorkDir) Root() string {
	return wd.root
}

func (wd *WorkDir) CreateFile(path string) error {
	fullPath := filepath.Join(wd.root, path)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	return file.Close()
}

func (wd *WorkDir) CreateDir(path string) error {
	fullPath := filepath.Join(wd.root, path)
	return os.MkdirAll(fullPath, 0755)
}

func (wd *WorkDir) WriteToFile(path string, content string) error {
	fullPath := filepath.Join(wd.root, path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (wd *WorkDir) AppendToFile(path string, content string) error {
	fullPath := filepath.Join(wd.root, path)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func (wd *WorkDir) ReadFile(path string) (string, error) {
	fullPath := filepath.Join(wd.root, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (wd *WorkDir) DeleteFile(path string) error {
	fullPath := filepath.Join(wd.root, path)
	return os.RemoveAll(fullPath)
}

func (wd *WorkDir) ListFilesRoot() []string {
	return wd.ListFilesRec()
}

func (wd *WorkDir) ListFilesRec() []string {
	var files []string
	filepath.WalkDir(wd.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(wd.root, path)
			if err != nil {
				return err
			}
			if strings.HasPrefix(rel, ".vc/") || strings.HasPrefix(rel, ".vc\\") {
				return nil
			}
			files = append(files, rel)
		}
		return nil
	})
	return files
}

func (wd *WorkDir) ListFilesIn(dir string) ([]string, error) {
	full := filepath.Join(wd.root, dir)
	var files []string
	err := filepath.WalkDir(full, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(wd.root, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func (wd *WorkDir) CatFile(path string) (string, error) {
	return wd.ReadFile(path)
}

func (wd *WorkDir) CopyFile(src, dst string) error {
	srcPath := filepath.Join(wd.root, src)
	dstPath := filepath.Join(wd.root, dst)
	return copy.Copy(srcPath, dstPath)
}

func (wd *WorkDir) Clone() *WorkDir {
	temp, err := os.MkdirTemp("", "workdir_clone")
	if err != nil {
		panic(err)
	}
	err = copy.Copy(wd.root, temp)
	if err != nil {
		panic(err)
	}
	return &WorkDir{root: temp}
}
