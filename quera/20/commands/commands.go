package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vc/workdir"
)

type VC struct {
	workdir     *workdir.WorkDir
	stagedFiles map[string]struct{}
	commits     []string
}

func (vc *VC) GetWorkDir() *workdir.WorkDir {
	return vc.workdir
}

func Init(wd *workdir.WorkDir) *VC {
	vc := &VC{
		workdir:     wd,
		stagedFiles: make(map[string]struct{}),
		commits:     []string{},
	}

	vcDir := filepath.Join(wd.Root(), ".vc")
	os.MkdirAll(filepath.Join(vcDir, "staging"), 0755)
	os.MkdirAll(filepath.Join(vcDir, "commits"), 0755)
	logPath := filepath.Join(vcDir, "messages.log")
	os.WriteFile(logPath, []byte(""), 0644)

	return vc
}

func (vc *VC) Add(paths ...string) error {
	for _, p := range paths {
		src := filepath.Join(vc.workdir.Root(), p)
		dst := filepath.Join(vc.workdir.Root(), ".vc", "staging", p)

		err := vc.workdir.CopyFile(src, dst)
		if err != nil {
			return fmt.Errorf("failed to add %s: %w", p, err)
		}
		vc.stagedFiles[p] = struct{}{}
	}
	return nil
}

func (vc *VC) AddAll() error {
	files := vc.workdir.ListFilesRec()
	return vc.Add(files...)
}

func (vc *VC) Commit(message string) error {
	stagingDir := filepath.Join(vc.workdir.Root(), ".vc", "staging")
	commitsDir := filepath.Join(vc.workdir.Root(), ".vc", "commits")

	if len(vc.stagedFiles) == 0 {
		return errors.New("nothing to commit")
	}

	commitID := len(vc.commits)
	newCommitPath := filepath.Join(commitsDir, strconv.Itoa(commitID))

	err := vc.workdir.CopyFile(stagingDir, newCommitPath)
	if err != nil {
		return err
	}

	logPath := filepath.Join(vc.workdir.Root(), ".vc", "messages.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(message + "\n")
	if err != nil {
		return err
	}

	vc.commits = append(vc.commits, message)
	vc.stagedFiles = make(map[string]struct{})

	err = os.RemoveAll(stagingDir)
	if err != nil {
		return err
	}
	return os.MkdirAll(stagingDir, 0755)
}

func (vc *VC) Log() []string {
	logPath := filepath.Join(vc.workdir.Root(), ".vc", "messages.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}

func (vc *VC) Checkout(id string) (*workdir.WorkDir, error) {
	num := 0
	for _, c := range id {
		if c == '~' || c == '^' {
			num++
		}
	}
	commitPath := filepath.Join(vc.workdir.Root(), ".vc", "commits", strconv.Itoa(num))
	info, err := os.Stat(commitPath)
	if err != nil || !info.IsDir() {
		return nil, errors.New("version not found")
	}

	tempDir, err := os.MkdirTemp("", "checkout")
	if err != nil {
		return nil, err
	}

	err = vc.workdir.CopyFile(commitPath, tempDir)
	if err != nil {
		return nil, err
	}

	return workdir.NewWorkDir(tempDir), nil
}

type Status struct {
	ModifiedFiles []string
	StagedFiles   []string
}

func (vc *VC) Status() Status {
	var modified []string
	var staged []string

	allFiles := vc.workdir.ListFilesRec()

	for _, f := range allFiles {
		stagingPath := filepath.Join(vc.workdir.Root(), ".vc", "staging", f)
		lastCommitNum := len(vc.commits) - 1
		commitPath := ""
		if lastCommitNum >= 0 {
			commitPath = filepath.Join(vc.workdir.Root(), ".vc", "commits", strconv.Itoa(lastCommitNum), f)
		}

		workContent, _ := vc.workdir.ReadFile(f)
		stagingContent, _ := vc.workdir.ReadFile(stagingPath)
		commitContent := ""
		if commitPath != "" {
			commitContent, _ = vc.workdir.ReadFile(commitPath)
		}

		if _, ok := vc.stagedFiles[f]; ok {
			if workContent != stagingContent {
				staged = append(staged, f)
			}
		} else {
			if workContent != commitContent {
				modified = append(modified, f)
			}
		}
	}

	return Status{ModifiedFiles: modified, StagedFiles: staged}
}
