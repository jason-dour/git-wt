// git-wt
//
// Helper for Git to manage worktrees as part of development flow.

package main

import (
	"os"

	"github.com/jason-dour/git-wt/internal/cobra/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
