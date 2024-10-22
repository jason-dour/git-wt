package git

import (
	"path/filepath"

	git "github.com/gogs/git-module"
	"github.com/jason-dour/git-wt/internal/cmn"
)

func ListWorktrees() ([]byte, error) {
	funcName := "git.ListWorktrees"
	cmn.Debug("%s: begin\n", funcName)

	cmd := git.NewCommand("worktree")
	cmd.AddArgs("list")
	cmn.Debug("%s: command: %s\n", funcName, cmd.String())

	output, err := cmd.RunInDir(filepath.Join(cmn.InitialDir, cmn.DefaultBranch))
	if err != nil {
		return nil, err
	}
	cmn.Debug("%s: output length: %d\n", funcName, len(output))

	cmn.Debug("%s: end\n", funcName)
	return output, nil
}
