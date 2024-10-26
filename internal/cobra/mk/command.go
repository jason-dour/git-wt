// Package mk implements the mk subcommand for git-wt.
package mk

import (
	"fmt"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command            = "mk"         // Command name.
	config  *cmn.CfgMk = &cmn.CfgMk{} // Configuration for the command.
	Cmd                = &cobra.Command{
		Use:     command + " worktree_name commit-ish",
		Short:   "Add a worktree to the project.",
		Long:    cmn.Basename + " " + command + " - Add a worktree to the project.",
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"make"},
		RunE:    run,
	} // Cobra command definition for the 'mk' command.
)

// init performs initialization for the 'mk' command.
func init() {
	Cmd.PersistentFlags().BoolVar(&config.Track, "track", false, "set up tracking mode")
	Cmd.PersistentFlags().StringVarP(&config.Branch, "b", "b", "", "create a new branch")
	Cmd.PersistentFlags().StringVarP(&config.BranchReset, "B", "B", "", "create or reset a branch")
	Cmd.PersistentFlags().BoolVar(&config.CheckoutNo, "no-checkout", false, "do not populate the new worktree")
	Cmd.PersistentFlags().BoolVarP(&config.Force, "force", "f", false, "checkout <branch> even if already checked out in other worktree")
	Cmd.PersistentFlags().BoolVarP(&config.Quiet, "quiet", "q", false, "suppress progress reporting")
}

// checkConfig scans config for proper use of flags.
func checkConfig() error {
	funcName := "checkConfig"
	cmn.Debug("%s: %s: begin", command, funcName)

	cmn.Debug("%s: %s: check mutually exclusive branch flags", command, funcName)
	if len(config.Branch) > 0 && len(config.BranchReset) > 0 {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("config: set branch with either -b or -B; don't use both")
	}

	cmn.Debug("%s: %s: check track has a new branch in flags", command, funcName)
	if config.Track && !(len(config.Branch) > 0 || len(config.BranchReset) > 0) {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("config: track requires new branch via -b or -B")
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}

// run is the main function for the 'mk' command.
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

	// Set the worktree name.
	wtName := args[0]
	cmn.Debug("%s: %s: worktree name: %s", command, funcName, wtName)

	// Set the commit-ish to be used.
	commitish := args[1]
	cmn.Debug("%s: %s: commit-ish: %s", command, funcName, commitish)

	// Check configuration.
	err = checkConfig()
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return err
	}

	// Add the worktree.
	output, err := git.WorktreeAdd(config, wtName, commitish)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return err
	}
	fmt.Print(string(output))

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}
