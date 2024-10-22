package git

import (
	"fmt"
	"os"

	git "github.com/gogs/git-module"
	"github.com/jason-dour/git-wt/internal/cmn"
)

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
