package commands

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"vc/workdir"

	"github.com/otiai10/copy"
)

type Status struct {
	ModifiedFiles []string
	StagedFiles   []string
}

type VC struct {
	wd *workdir.WorkDir
}

func Init(wd *workdir.WorkDir) *VC {
	vcDir := filepath.Join(wd.Root(), ".vc")
	os.MkdirAll(filepath.Join(vcDir, "staging"), 0755)
	os.MkdirAll(filepath.Join(vcDir, "commits"), 0755)
	wd.WriteToFile(".vc/messages.log", "")
	return &VC{wd: wd}
}

func (vc *VC) GetWorkDir() *workdir.WorkDir {
	return vc.wd
}

func (vc *VC) Add(paths ...string) error {
	for _, p := range paths {
		dst := filepath.Join(".vc", "staging", p)
		err := vc.wd.CreateDir(filepath.Dir(dst))
		if err != nil {
			return err
		}
		err = vc.wd.CopyFile(p, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vc *VC) AddAll() error {
	files := vc.wd.ListFilesRec()
	var filtered []string
	for _, f := range files {
		if !strings.HasPrefix(f, ".vc/") && !strings.HasPrefix(f, ".vc\\") {
			filtered = append(filtered, f)
		}
	}
	return vc.Add(filtered...)
}

func (vc *VC) Commit(message string) error {
	stagingPath := filepath.Join(vc.wd.Root(), ".vc", "staging")
	entries, err := os.ReadDir(stagingPath)
	if err != nil || len(entries) == 0 {
		return errors.New("nothing to commit")
	}

	commitsDir := filepath.Join(vc.wd.Root(), ".vc", "commits")
	commitEntries, err := os.ReadDir(commitsDir)
	if err != nil {
		return err
	}
	next := len(commitEntries)
	dst := filepath.Join(vc.wd.Root(), ".vc", "commits", strconv.Itoa(next))

	err = copy.Copy(stagingPath, dst)
	if err != nil {
		return err
	}

	content, err := vc.wd.ReadFile(".vc/messages.log")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content += message + "\n"
	if err := vc.wd.WriteToFile(".vc/messages.log", content); err != nil {
		if os.IsNotExist(err) {
			if err := vc.wd.CreateFile(".vc/messages.log"); err != nil {
				return err
			}
			if err := vc.wd.WriteToFile(".vc/messages.log", content); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := os.RemoveAll(stagingPath); err != nil {
		return err
	}
	if err := os.MkdirAll(stagingPath, 0755); err != nil {
		return err
	}

	return nil
}

func (vc *VC) Log() []string {
	content, err := vc.wd.ReadFile(".vc/messages.log")
	if err != nil {
		return []string{}
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return []string{}
	}
	lines := strings.Split(content, "\n")
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (vc *VC) Status() Status {
	stagingDir := filepath.Join(vc.wd.Root(), ".vc", "staging")
	stagedFullPaths, _ := vc.wd.ListFilesIn(stagingDir)
	var staged []string
	for _, s := range stagedFullPaths {
		rel, err := filepath.Rel(stagingDir, filepath.Join(vc.wd.Root(), s))
		if err == nil {
			staged = append(staged, rel)
		}
	}

	modified := []string{}
	commitsDir := filepath.Join(vc.wd.Root(), ".vc", "commits")
	commitEntries, err := os.ReadDir(commitsDir)
	if err != nil || len(commitEntries) == 0 {
		return Status{StagedFiles: staged, ModifiedFiles: modified}
	}

	latest := len(commitEntries) - 1
	base := filepath.Join(vc.wd.Root(), ".vc", "commits", strconv.Itoa(latest))
	commitFiles, _ := vc.wd.ListFilesIn(base)

	for _, f := range commitFiles {
		rel, err := filepath.Rel(base, f)
		if err != nil {
			continue
		}

		currentPath := filepath.Join(vc.wd.Root(), rel)
		currentBytes, err1 := os.ReadFile(currentPath)
		oldBytes, err2 := os.ReadFile(f)

		if err1 == nil && err2 == nil && !bytes.Equal(currentBytes, oldBytes) {
			if !contains(staged, rel) {
				modified = append(modified, rel)
			}
		}
	}

	allFiles := vc.wd.ListFilesRec()
	for _, f := range allFiles {
		fFull := filepath.Join(vc.wd.Root(), f)
		isInCommit := false
		for _, cf := range commitFiles {
			if cf == fFull {
				isInCommit = true
				break
			}
		}
		if !isInCommit && !contains(staged, f) {
			modified = append(modified, f)
		}
	}

	return Status{StagedFiles: staged, ModifiedFiles: modified}
}

func (vc *VC) Checkout(id string) (*workdir.WorkDir, error) {
	count := strings.Count(id, "~") + strings.Count(id, "^")
	commitsDir := filepath.Join(vc.wd.Root(), ".vc", "commits")
	entries, err := os.ReadDir(commitsDir)
	if err != nil {
		return nil, errors.New("no commits found")
	}

	target := len(entries) - 1 - count
	if target < 0 || target >= len(entries) {
		return nil, errors.New("version not found")
	}

	src := filepath.Join(vc.wd.Root(), ".vc", "commits", strconv.Itoa(target))
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil, errors.New("commit directory not found")
	}

	tmpDir, err := os.MkdirTemp("", "checkout")
	if err != nil {
		return nil, err
	}

	// استفاده از copy.Copy برای کپی دایرکتوری
	err = copy.Copy(src, tmpDir)
	if err != nil {
		return nil, err
	}

	return workdir.NewWorkDir(tmpDir), nil
}
