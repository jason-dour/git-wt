// Package rm implements the rm subcommand for git-wt.
package rm

import (
	"fmt"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command            = "rm"         // Command name.
	config  *cmn.CfgRm = &cmn.CfgRm{} // Configuration for the command.
	Cmd                = &cobra.Command{
		Use:     command,
		Short:   "Remove a worktree from the project.",
		Long:    cmn.Basename + " " + command + " - Remove a worktree from the project.",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"remove", "del", "delete"},
		RunE:    run,
	} // Cobra command definition for the 'rm' command.
)

// init performs initialization for the 'rm' command.
func init() {
	Cmd.PersistentFlags().BoolVarP(&config.Force, "force", "f", false, "force removal even if worktree is dirty or locked")
}

// run provides the core execution of the 'rm' command.
func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin\n", command, funcName)
	cmn.Debug("%s: %s: config: %#v\n", command, funcName, config)
	cmn.Debug("%s: %s: args: %v\n", command, funcName, args)

	// Set the worktree name.
	wtName := args[0]
	cmn.Debug("%s: %s: worktree name: %s\n", command, funcName, wtName)

	// Remove the worktree.
	output, err := git.WorktreeRemove(config, wtName)
	if err != nil {
		return err
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}