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
	cmn.Debug("%s: begin", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module: ")
	}

	err := git.Clone(url, path, git.CloneOptions{
		Branch: branch,
	})
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return fmt.Errorf("could not clone repo: %s", err.Error())
	}

	cmn.Debug("%s: end", funcName)
	return nil
}

// DeleteBranch will delete a local branch from the repository.
func DeleteBranch(branch string) error {
	funcName := "git.DeleteBranch"
	cmn.Debug("%s: begin", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module: ")
	}

	err := git.DeleteBranch(
		filepath.Join(cmn.Config.ProjectDir, cmn.Config.DefaultBranch),
		branch,
		git.DeleteBranchOptions{
			Force: true,
		})
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return fmt.Errorf("error deleting branch: %s", err.Error())
	}

	cmn.Debug("%s: end", funcName)
	return nil
}

// GetBranches will retrieve local branches from the repository.
func GetBranches() ([]string, error) {
	funcName := "git.GetBranches"
	cmn.Debug("%s: begin", funcName)

	repository, err := git.Open(filepath.Join(cmn.Config.ProjectDir, cmn.Config.DefaultBranch))
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return []string{}, fmt.Errorf("error opening repository: %s", err.Error())
	}

	branches, err := repository.Branches()
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return []string{}, fmt.Errorf("error getting branches: %s", err.Error())
	}
	cmn.Debug("%s: branches: %v", funcName, branches)

	cmn.Debug("%s: end", funcName)
	return branches, nil
}

// GetRemote will get the remote URL for the origin.
func GetRemote(directory string) (string, error) {
	funcName := "git.getRemote"
	cmn.Debug("%s: begin", funcName)

	// g remote get-url origin

	cmd := git.NewCommand("remote")
	cmd.AddArgs("get-url", "origin")

	cmn.Debug("%s: command: %s", funcName, cmd.String())

	output, err := cmd.RunInDir(directory)
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return "", err
	}
	cmn.Debug("%s: output length: %d", funcName, len(output))

	remote := ""
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		cmn.Debug("%s: output line: %s", funcName, scanner.Text())
		remote = scanner.Text()
		break
	}
	cmn.Debug("%s: remote: %s", funcName, remote)

	cmn.Debug("%s: end", funcName)
	return remote, nil
}

// GetRemoteDefaultBranch will retrieve the default branch from a remote git repository.
func GetRemoteDefaultBranch(url string) (string, error) {
	funcName := "git.GetRemoteDefaultBranch"
	cmn.Debug("%s: begin", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module: ")
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
		cmn.Debug("%s: error: end", funcName)
		return "", fmt.Errorf("could not retrieve HEAD ref from remote: %s", err.Error())
	}

	// Iterate through returned refs.
	for i, v := range refs {
		cmn.Debug("%s: refs[%v]: { id: %v, refspec: %v }", funcName, i, v.ID, v.Refspec)
		if v.ID == "ref:" {
			defaultBranch = filepath.Base(v.Refspec)
			cmn.Debug("%s: found default branch: %s", funcName, defaultBranch)
		}
	}
	cmn.Debug("%s: defaultBranch: %v", funcName, defaultBranch)

	// If we couldn't find the default branch, return error.
	if defaultBranch == "" {
		cmn.Debug("%s: error: end", funcName)
		return "", fmt.Errorf("could not determine default branch from HEAD ref")
	}

	// Return default branch.
	cmn.Debug("%s: end", funcName)
	return defaultBranch, nil
}

// GetRemoteRefId will retrieve the commit id for a given ref.
func GetRemoteRefId(url string, ref string) (string, error) {
	funcName := "git.GetRemoteRefId"
	cmn.Debug("%s: begin", funcName)

	// If Debug, set debug for git-module.
	if cmn.Config.DebugFlag {
		git.SetOutput(os.Stderr)
		git.SetPrefix("debug: git-module: ")
	}

	// Run ls-remote to grab all head refs.
	refs, err := git.LsRemote(url, git.LsRemoteOptions{
		Patterns: []string{
			"refs/" + ref,
		},
	})
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return "", fmt.Errorf("could not retrieve HEAD ref from remote: %s", err.Error())
	}

	if len(refs) > 1 {
		cmn.Debug("%s: error: end", funcName)
		return "", fmt.Errorf("got more than one matching ref from remote: %d", len(refs))
	}

	cmn.Debug("%s: found ref: %#v", funcName, refs[0])

	cmn.Debug("%s: end", funcName)
	return refs[0].ID, nil
}

// WorktreeAdd will add a worktree to the project.
func WorktreeAdd(config *cmn.CfgMk, worktree string, commitish string) ([]byte, error) {
	funcName := "git.WorktreeAdd"
	cmn.Debug("%s: begin", funcName)
	cmn.Debug("%s: config: %#v", funcName, config)

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
	if len(config.RefId) > 0 {
		cmd.AddArgs(config.RefId)
	} else {
		cmd.AddArgs(commitish)
	}

	cmn.Debug("%s: command: %s", funcName, cmd.String())

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return nil, err
	}
	cmn.Debug("%s: output length: %d", funcName, len(output))

	cmn.Debug("%s: end", funcName)
	return output, nil
}

// WorktreeList will list all worktrees in the project.
func WorktreeList(porcelain bool) ([]byte, error) {
	funcName := "git.WorktreeList"
	cmn.Debug("%s: begin", funcName)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("list")

	if porcelain {
		cmd.AddArgs("--porcelain")
	}

	cmn.Debug("%s: command: %s", funcName, cmd.String())

	runDir := ""
	if cmn.Config.InitialDir == cmn.Config.ProjectDir {
		runDir = filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch)
	} else {
		runDir = cmn.Config.InitialDir
	}
	output, err := cmd.RunInDir(runDir)
	if err != nil {
		cmn.Debug("%s: error: end", funcName)
		return nil, err
	}
	cmn.Debug("%s: output length: %d", funcName, len(output))

	cmn.Debug("%s: end", funcName)
	return output, nil
}

// WorktreeMove will move a worktree within the project.
func WorktreeMove(config *cmn.CfgMv, wtOriginal string, wtNew string) ([]byte, error) {
	funcName := "git.WorktreeMove"
	cmn.Debug("%s: begin", funcName)
	cmn.Debug("%s: config: %#v", funcName, config)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("move")

	if config.Force {
		cmd.AddArgs("--force")
	}

	cmd.AddArgs(filepath.Join(cmn.Config.ProjectDir, wtOriginal), filepath.Join(cmn.Config.ProjectDir, wtNew))

	cmn.Debug("%s: command: %s", funcName, cmd.String())

	if strings.HasPrefix(cmn.Config.InitialDir, filepath.Join(cmn.Config.ProjectDir, wtOriginal)) {
		cmn.Debug("%s: error: end", funcName)
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
		cmn.Debug("%s: error: end", funcName)
		return nil, err
	}
	cmn.Debug("%s: output length: %d", funcName, len(output))

	cmn.Debug("%s: end", funcName)
	return output, nil
}

// WorktreeRemove will remove a worktree from the project.
func WorktreeRemove(config *cmn.CfgRm, worktree string) ([]byte, error) {
	funcName := "git.WorktreeRemove"
	cmn.Debug("%s: begin", funcName)
	cmn.Debug("%s: config: %#v", funcName, config)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("remove")

	if config.Force {
		cmd.AddArgs("--force")
	}

	cmd.AddArgs(worktree)

	cmn.Debug("%s: command: %s", funcName, cmd.String())

	if strings.HasPrefix(cmn.Config.InitialDir, filepath.Join(cmn.Config.ProjectDir, worktree)) {
		cmn.Debug("%s: error: end", funcName)
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
		cmn.Debug("%s: error: end", funcName)
		return nil, err
	}
	cmn.Debug("%s: output length: %d", funcName, len(output))

	cmn.Debug("%s: end", funcName)
	return output, nil
}
