// git-wt is a helper for Git to manage worktrees as part of a development flow.
//
// In addition to providing many of the 'git worktree' command set but with
// shorter syntax, it also provides a clone feature that structures your clone
// to readily make use of worktrees in one project directory.
//
// Usage:
//
// git-wt [command]
//
// Available Commands:
//
//	clone       Clone a repo for a git-wt workflow.
//	completion  Generate the autocompletion script for the specified shell
//	help        Help about any command
//	ls          List worktrees for the project.
//
// Flags:
//
//	-d, --debug     enable debug mode
//	-h, --help      help for git-wt
//	-v, --version   version for git-wt
//
// Use "git-wt [command] --help" for more information about a command.
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
