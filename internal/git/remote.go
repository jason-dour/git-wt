package git

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/gogs/git-module"
	"github.com/jason-dour/git-wt/internal/cmn"
)

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
