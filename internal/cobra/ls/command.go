// Package list implements the list subcommand for git-wt.
package ls

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
	} // Cobra command definition for the 'ls' command.
)

// run is the main function for the 'ls' command.
func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin", command, funcName)

	// Load global configuration.
	cmn.Debug("%s: %s: loading global config", command, funcName)
	err := cmn.InitConfig()
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error loading configuration: %s", err.Error())
	}
	cmn.Debug("%s: %s: global config: %#v", command, funcName, cmn.Config)

	output, err := git.WorktreeList(false)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error listing worktrees: %s", err.Error())
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}
