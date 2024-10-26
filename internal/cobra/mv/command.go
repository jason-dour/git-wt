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

// run provides the core execution of the 'mv' command.
func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	// Load global configuration.
	cmn.Debug("%s: %s: loading global config\n", command, funcName)
	err := cmn.InitConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %s", err.Error())
	}
	cmn.Debug("%s: %s: global config: %#v\n", command, funcName, cmn.Config)

	cmn.Debug("%s: %s: config: %#v\n", command, funcName, config)
	cmn.Debug("%s: %s: args: %v\n", command, funcName, args)

	// Set the worktree current name.
	wtCurr := args[0]
	cmn.Debug("%s: %s: worktree name: %s\n", command, funcName, wtCurr)

	// Set the worktree current name.
	wtNew := args[1]
	cmn.Debug("%s: %s: new worktree name: %s\n", command, funcName, wtNew)

	// Remove the worktree.
	output, err := git.WorktreeMove(config, wtCurr, wtNew)
	if err != nil {
		return err
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}
