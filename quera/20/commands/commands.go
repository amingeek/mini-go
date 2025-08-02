// commands/commands.go
package commands

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"vc/workdir"
)

type Commit struct {
	Message string
	Files   map[string]string
}

type VC struct {
	mu          sync.Mutex
	workDir     *workdir.WorkDir
	commits     []*Commit
	stagingArea map[string]string
}

// Init initializes a new version control instance on given WorkDir
func Init(wd *workdir.WorkDir) *VC {
	return &VC{
		workDir:     wd,
		commits:     []*Commit{},
		stagingArea: make(map[string]string),
	}
}

// GetWorkDir returns underlying WorkDir
func (vc *VC) GetWorkDir() *workdir.WorkDir {
	return vc.workDir
}

// Add stages specified files for next commit
func (vc *VC) Add(files ...string) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	for _, f := range files {
		content, err := vc.workDir.CatFile(f)
		if err != nil {
			return fmt.Errorf("file %s not found: %w", f, err)
		}
		vc.stagingArea[f] = content
	}
	return nil
}

// AddAll stages all existing files
func (vc *VC) AddAll() error {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	files := vc.workDir.ListFilesRoot()
	for _, f := range files {
		content, err := vc.workDir.CatFile(f)
		if err != nil {
			return err
		}
		vc.stagingArea[f] = content
	}
	return nil
}

// Commit creates a new commit with given message. Always records a new commit, even if no changes staged.
func (vc *VC) Commit(message string) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	// Build snapshot from last commit
	commitFiles := map[string]string{}
	if len(vc.commits) > 0 {
		for f, c := range vc.commits[len(vc.commits)-1].Files {
			commitFiles[f] = c
		}
	}
	// Overlay staged changes
	for f, c := range vc.stagingArea {
		commitFiles[f] = c
	}
	// Remove files deleted in workdir
	for f := range commitFiles {
		if _, err := vc.workDir.CatFile(f); err != nil {
			delete(commitFiles, f)
		}
	}
	// Append commit
	newCommit := &Commit{Message: message, Files: commitFiles}
	vc.commits = append(vc.commits, newCommit)
	// Clear staging
	vc.stagingArea = map[string]string{}
	return nil
}

// StatusResult holds lists of modified and staged files
type StatusResult struct {
	ModifiedFiles []string
	StagedFiles   []string
}

// Status returns current status: which files are modified and which are staged
func (vc *VC) Status() StatusResult {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	modified := []string{}
	staged := []string{}
	// If no commit, everything is clean
	if len(vc.commits) == 0 {
		return StatusResult{ModifiedFiles: modified, StagedFiles: staged}
	}
	lastFiles := vc.commits[len(vc.commits)-1].Files
	// Collect all file paths
	all := map[string]struct{}{}
	for _, f := range vc.workDir.ListFilesRoot() {
		all[f] = struct{}{}
	}
	for f := range lastFiles {
		all[f] = struct{}{}
	}
	for f := range vc.stagingArea {
		all[f] = struct{}{}
	}
	// Evaluate each file
	for f := range all {
		workContent, workErr := vc.workDir.CatFile(f)
		lastContent, hasLast := lastFiles[f]
		stagedContent, hasStage := vc.stagingArea[f]
		inWork := workErr == nil
		inLast := hasLast
		inStage := hasStage
		if inStage {
			// always list in staged
			staged = append(staged, f)
			// if changed after staging, also list modified
			if inWork && workContent != stagedContent {
				modified = append(modified, f)
			}
		} else {
			// not staged
			if inLast && inWork && workContent != lastContent {
				modified = append(modified, f)
			} else if !inLast && inWork {
				modified = append(modified, f)
			}
		}
	}
	sort.Strings(modified)
	sort.Strings(staged)
	return StatusResult{ModifiedFiles: modified, StagedFiles: staged}
}

// Log returns commit messages in reverse chronological order
func (vc *VC) Log() []string {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	res := []string{}
	for i := len(vc.commits) - 1; i >= 0; i-- {
		res = append(res, vc.commits[i].Message)
	}
	return res
}

// Checkout returns a cloned WorkDir at given revision (~N or ^ chains)
func (vc *VC) Checkout(rev string) (*workdir.WorkDir, error) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	n := 0
	switch {
	case rev == "^":
		n = 1
	case rev == "^^":
		n = 2
	case rev == "^^^":
		n = 3
	case strings.HasPrefix(rev, "~"):
		if _, err := fmt.Sscanf(rev, "~%d", &n); err != nil {
			return nil, errors.New("invalid revision format")
		}
	case rev == "":
		n = 0
	default:
		return nil, errors.New("unsupported revision format")
	}
	if n >= len(vc.commits) {
		return nil, errors.New("revision out of range")
	}
	commit := vc.commits[len(vc.commits)-1-n]
	newWD := workdir.InitEmptyWorkDir()
	for f, content := range commit.Files {
		if err := newWD.CreateFile(f); err != nil {
			return nil, err
		}
		if err := newWD.WriteToFile(f, content); err != nil {
			return nil, err
		}
	}
	return newWD, nil
}
