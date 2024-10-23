// Package mk implements the mk subcommand for git-wt.
package mk

import (
	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/spf13/cobra"
)

var (
	command = "mk" // Command name.
	Cmd     = &cobra.Command{
		Use:     command + " repo_url",
		Short:   "Add a worktree to the project.",
		Long:    cmn.Basename + " " + command + " - Add a worktree to the project.",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"make"},
		RunE:    run,
	} // Cobra command definition for the 'clone' command.
)

func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin\n", command, funcName)
	cmn.Debug("%s: %s: args: %v\n", command, funcName, args)

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}
