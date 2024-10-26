// Package mv implements the mv subcommand for git-wt.
package mv

import (
	"fmt"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command            = "mv"         // Command name.
	config  *cmn.CfgMv = &cmn.CfgMv{} // Configuration for the command.
	Cmd                = &cobra.Command{
		Use:     command + " worktree-name new-worktree-name",
		Short:   "Move a worktree within the project.",
		Long:    cmn.Basename + " " + command + " - Move a worktree within the project.",
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"move", "ren", "rename"},
		RunE:    run,
	} // Cobra command definition for the 'mv' command.
)

// init performs initialization for the 'mv' command.
func init() {
	Cmd.PersistentFlags().BoolVarP(&config.Force, "force", "f", false, "force move even if worktree is dirty or locked")
}

// run is the main function for the 'mv' command.
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

	cmn.Debug("%s: %s: config: %#v", command, funcName, config)
	cmn.Debug("%s: %s: args: %v", command, funcName, args)

	// Set the worktree current name.
	wtCurr := args[0]
	cmn.Debug("%s: %s: worktree name: %s", command, funcName, wtCurr)

	// Set the worktree current name.
	wtNew := args[1]
	cmn.Debug("%s: %s: new worktree name: %s", command, funcName, wtNew)

	// Remove the worktree.
	output, err := git.WorktreeMove(config, wtCurr, wtNew)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return err
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}
