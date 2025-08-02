// workdir/workdir.go
package workdir

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type WorkDir struct {
	root string
}

func InitEmptyWorkDir() *WorkDir {
	dir, err := ioutil.TempDir("", "workdir")
	if err != nil {
		panic(err)
	}
	return &WorkDir{root: dir}
}

func (wd *WorkDir) resolve(path string) string {
	return filepath.Join(wd.root, filepath.FromSlash(path))
}

func (wd *WorkDir) CreateFile(path string) error {
	full := wd.resolve(path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	f, err := os.Create(full)
	if err != nil {
		return err
	}
	return f.Close()
}

func (wd *WorkDir) CreateDir(path string) error {
	return os.MkdirAll(wd.resolve(path), 0755)
}

func (wd *WorkDir) WriteToFile(path, content string) error {
	full := wd.resolve(path)
	if _, err := os.Stat(full); errors.Is(err, os.ErrNotExist) {
		return errors.New("file does not exist")
	}
	return ioutil.WriteFile(full, []byte(content), 0644)
}

func (wd *WorkDir) AppendToFile(path, content string) error {
	full := wd.resolve(path)
	if _, err := os.Stat(full); errors.Is(err, os.ErrNotExist) {
		return errors.New("file does not exist")
	}
	f, err := os.OpenFile(full, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func (wd *WorkDir) CatFile(path string) (string, error) {
	full := wd.resolve(path)
	data, err := ioutil.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (wd *WorkDir) ListFilesRoot() []string {
	var files []string
	filepath.Walk(wd.root, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		r, err := filepath.Rel(wd.root, p)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(r))
		return nil
	})
	return files
}

func (wd *WorkDir) ListFilesIn(dir string) ([]string, error) {
	var files []string
	full := wd.resolve(dir)
	err := filepath.Walk(full, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		r, err := filepath.Rel(wd.root, p)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(r))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (wd *WorkDir) Clone() *WorkDir {
	newWD := InitEmptyWorkDir()
	copyDir(wd.root, newWD.root)
	return newWD
}

func copyDir(src, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			copyDir(srcPath, dstPath)
		} else {
			data, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
