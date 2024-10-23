// Package list implements the list subcommand for git-wt.
package list

import (
	"fmt"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command = "ls" // Command name.
	Cmd     = &cobra.Command{
		Use:     command,
		Short:   "List worktrees for the project.",
		Long:    cmn.Basename + " " + command + " - List worktrees for the project.",
		Args:    cobra.NoArgs,
		Aliases: []string{"list"},
		RunE:    run,
	} // Cobra command definition for the 'list' command.
)

func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	output, err := git.WorktreeList()
	if err != nil {
		cmn.Exit(20, "error listing worktrees: %s", err.Error())
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}
