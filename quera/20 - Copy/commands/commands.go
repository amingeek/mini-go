// commands.go
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
	wd.CreateFile(".vc/messages.log")

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

func getNextCommitID(commitsDir string) int {
	entries, err := os.ReadDir(commitsDir)
	if err != nil {
		return 0
	}
	max := -1
	for _, e := range entries {
		if e.IsDir() {
			if id, err := strconv.Atoi(e.Name()); err == nil && id > max {
				max = id
			}
		}
	}
	return max + 1
}

func (vc *VC) Commit(message string) error {
	stagingPath := filepath.Join(vc.wd.Root(), ".vc", "staging")

	hasFile := false
	filepath.Walk(stagingPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hasFile = true
		}
		return nil
	})
	if !hasFile {
		return errors.New("nothing to commit")
	}

	commitsDir := filepath.Join(vc.wd.Root(), ".vc", "commits")
	next := getNextCommitID(commitsDir)
	dst := filepath.Join(commitsDir, strconv.Itoa(next))

	err := copy.Copy(stagingPath, dst)
	if err != nil {
		return err
	}

	err = vc.wd.AppendToFile(".vc/messages.log", message+"\n")
	if err != nil {
		if os.IsNotExist(err) {
			if err := vc.wd.CreateFile(".vc/messages.log"); err != nil {
				return err
			}
			err = vc.wd.AppendToFile(".vc/messages.log", message+"\n")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := os.RemoveAll(stagingPath); err != nil {
		return err
	}
	return os.MkdirAll(stagingPath, 0755)
}

func (vc *VC) Log() []string {
	content, err := vc.wd.ReadFile(".vc/messages.log")
	if err != nil || strings.TrimSpace(content) == "" {
		return []string{}
	}
	lines := strings.Split(strings.TrimSpace(content), "\n")
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
	stagedPaths, _ := vc.wd.ListFilesIn(".vc/staging")
	staged := append([]string{}, stagedPaths...)
	modified := []string{}

	commitsDir := filepath.Join(vc.wd.Root(), ".vc", "commits")
	entries, err := os.ReadDir(commitsDir)
	if err != nil || len(entries) == 0 {
		return Status{StagedFiles: staged, ModifiedFiles: []string{}}
	}

	latestCommit := filepath.Join(".vc/commits", strconv.Itoa(getNextCommitID(commitsDir)-1))
	commitFiles, _ := vc.wd.ListFilesIn(latestCommit)

	for _, relPath := range commitFiles {
		commitPath := filepath.Join(vc.wd.Root(), latestCommit, relPath)
		currentPath := filepath.Join(vc.wd.Root(), relPath)

		oldBytes, err1 := os.ReadFile(commitPath)
		newBytes, err2 := os.ReadFile(currentPath)

		if err1 == nil && err2 == nil && !bytes.Equal(oldBytes, newBytes) {
			if contains(staged, relPath) {
				modified = append(modified, relPath)
				newStaged := []string{}
				for _, s := range staged {
					if s != relPath {
						newStaged = append(newStaged, s)
					}
				}
				staged = newStaged
			} else {
				modified = append(modified, relPath)
			}
		}
	}

	allFiles := vc.wd.ListFilesRec()
	for _, f := range allFiles {
		if !contains(commitFiles, f) && !contains(staged, f) && !contains(modified, f) {
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

	numCommits := 0
	for _, entry := range entries {
		if entry.IsDir() {
			numCommits++
		}
	}

	target := numCommits - 1 - count
	if target < 0 {
		return nil, errors.New("version not found")
	}

	src := filepath.Join(commitsDir, strconv.Itoa(target))

	tmpDir, err := os.MkdirTemp("", "checkout")
	if err != nil {
		return nil, err
	}

	err = copy.Copy(src, tmpDir)
	if err != nil {
		return nil, err
	}

	return workdir.NewWorkDir(tmpDir), nil
}
