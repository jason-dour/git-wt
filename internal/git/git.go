// Package git implements abstraction to the module git-module. It also provides
// for worktree commands not implemented in git-module.
package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogs/git-module"
	"github.com/jason-dour/git-wt/internal/cmn"
)

// Clone will clone a git repository, checkout a branch, to a path provided.
func Clone(url string, branch string, path string) error {
	funcName := "git.Clone"
	cmn.Debug("%s: begin\n", funcName)

	// If Debug, set debug for git-module.
	if cmn.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module:")
	}

	err := git.Clone(url, path, git.CloneOptions{
		Branch: branch,
	})
	if err != nil {
		return fmt.Errorf("could not clone repo: %s", err.Error())
	}

	return nil
}

// GetDefaultBranch will retrieve the default branch from a remote git repository.
func GetDefaultBranch(url string) (string, error) {
	funcName := "git.GetDefaultBranch"
	cmn.Debug("%s: begin\n", funcName)

	// If Debug, set debug for git-module.
	if cmn.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module:")
	}

	// Initialize the return value.
	defaultBranch := ""

	// Run ls-remote to grab HEAD symref.
	refs, err := git.LsRemote(url, git.LsRemoteOptions{
		Patterns: []string{
			"HEAD",
		},
		CommandOptions: git.CommandOptions{
			Args: []string{
				"--symref",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("could not retrieve HEAD ref from remote: %s", err.Error())
	}

	// Iterate through returned refs.
	for i, v := range refs {
		cmn.Debug("%s: refs[%v]: { id: %v, refspec: %v }\n", funcName, i, v.ID, v.Refspec)
		if v.ID == "ref:" {
			defaultBranch = filepath.Base(v.Refspec)
			cmn.Debug("%s: found default branch: %s\n", funcName, defaultBranch)
		}
	}
	cmn.Debug("%s: defaultBranch: %v\n", funcName, defaultBranch)

	// If we couldn't find the default branch, return error.
	if defaultBranch == "" {
		return "", fmt.Errorf("could not determine default branch from HEAD ref")
	}

	// Return default branch.
	cmn.Debug("%s: end\n", funcName)
	return defaultBranch, nil
}

// WorktreeList will list all worktrees in the project.
func WorktreeList() ([]byte, error) {
	funcName := "git.ListWorktrees"
	cmn.Debug("%s: begin\n", funcName)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("list")
	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	runDir := ""
	if cmn.InitialDir == cmn.ProjectDir {
		runDir = filepath.Join(cmn.InitialDir, cmn.DefaultBranch)
	} else {
		runDir = cmn.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}
