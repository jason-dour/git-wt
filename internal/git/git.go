// Package git implements abstraction to the module git-module. It also provides
// for worktree commands not implemented in git-module.
package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogs/git-module"
	"github.com/jason-dour/git-wt/internal/cmn"
)

// Clone will clone a git repository, checkout a branch, to a path provided.
func Clone(url string, branch string, path string) error {
	funcName := "git.Clone"
	cmn.Debug("%s: begin\n", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module: ")
	}

	err := git.Clone(url, path, git.CloneOptions{
		Branch: branch,
	})
	if err != nil {
		return fmt.Errorf("could not clone repo: %s", err.Error())
	}

	cmn.Debug("%s: end\n", funcName)
	return nil
}

// GetDefaultBranch will retrieve the default branch from a remote git repository.
func GetDefaultBranch(url string) (string, error) {
	funcName := "git.GetDefaultBranch"
	cmn.Debug("%s: begin\n", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
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

// GetRemote will get the remote URL for the origin.
func GetRemote(directory string) (string, error) {
	funcName := "git.getRemote"
	cmn.Debug("%s: begin\n", funcName)

	// g remote get-url origin

	cmd := git.NewCommand("remote")
	cmd.AddArgs("get-url", "origin")

	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	output, err := cmd.RunInDir(directory)
	if err != nil {
		return "", err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	remote := ""
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		cmn.Debug("%s: output line: %s\n", funcName, scanner.Text())
		remote = scanner.Text()
		break
	}
	cmn.Debug("%s: remote: %s\n", funcName, remote)

	cmn.Debug("%s: end\n", funcName)
	return remote, nil
}

// WorktreeAdd will add a worktree to the project.
func WorktreeAdd(config *cmn.CfgMk, worktree string, commitish string) ([]byte, error) {
	funcName := "git.WorktreeAdd"
	cmn.Debug("%s: begin\n", funcName)
	cmn.Debug("%s: config: %#v\n", funcName, config)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("add")

	if config.Quiet {
		cmd.AddArgs("--quiet")
	}
	if config.Force {
		cmd.AddArgs("--force")
	}
	if config.CheckoutNo {
		cmd.AddArgs("--no-checkout")
	}
	if config.Track {
		cmd.AddArgs("--track")
	}
	if len(config.Branch) > 0 {
		cmd.AddArgs("-b", config.Branch)
	}
	if len(config.BranchReset) > 0 {
		cmd.AddArgs("-B", config.BranchReset)
	}

	cmd.AddArgs(filepath.Join(cmn.Config.ProjectDir, worktree))
	cmd.AddArgs(commitish)

	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}

// WorktreeList will list all worktrees in the project.
func WorktreeList() ([]byte, error) {
	funcName := "git.WorktreeList"
	cmn.Debug("%s: begin\n", funcName)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("list")
	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}

// WorktreeMove will move a worktree within the project.
func WorktreeMove(config *cmn.CfgMv, wtOriginal string, wtNew string) ([]byte, error) {
	funcName := "git.WorktreeMove"
	cmn.Debug("%s: begin\n", funcName)
	cmn.Debug("%s: config: %#v\n", funcName, config)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("move")

	if config.Force {
		cmd.AddArgs("--force")
	}

	cmd.AddArgs(filepath.Join(cmn.Config.ProjectDir, wtOriginal), filepath.Join(cmn.Config.ProjectDir, wtNew))

	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	if strings.HasPrefix(cmn.Config.InitialDir, filepath.Join(cmn.Config.ProjectDir, wtOriginal)) {
		return nil, fmt.Errorf("cannot move worktree; current working directory within worktree")
	}

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}

// WorktreeRemove will remove a worktree from the project.
func WorktreeRemove(config *cmn.CfgRm, worktree string) ([]byte, error) {
	funcName := "git.WorktreeRemove"
	cmn.Debug("%s: begin\n", funcName)
	cmn.Debug("%s: config: %#v\n", funcName, config)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("remove")

	if config.Force {
		cmd.AddArgs("--force")
	}

	cmd.AddArgs(worktree)

	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	if strings.HasPrefix(cmn.Config.InitialDir, filepath.Join(cmn.Config.ProjectDir, worktree)) {
		return nil, fmt.Errorf("cannot remove worktree; current working directory within worktree")
	}

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}
