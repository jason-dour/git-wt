// Package cl implements the cl subcommand for git-wt.
package cl

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command = "cl" // Command name.
	Cmd     = &cobra.Command{
		Use:     command + " repo_url",
		Short:   "Clone a repo for a git-wt workflow.",
		Long:    cmn.Basename + " " + command + " - Clone a repo for a git-wt workflow.",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"clone"},
		RunE:    run,
	} // Cobra command definition for the 'clone' command.
)

// run is the main function for the 'cl' command.
func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin", command, funcName)
	cmn.Debug("%s: %s: args: %v", command, funcName, args)

	fmt.Printf("Cloning %s.\n", args[0])

	// Get default branch from remote repository.
	cmn.Debug("%s: %s: retrieving default branch from remote", command, funcName)
	defaultBranch, err := git.GetRemoteDefaultBranch(args[0])
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("could not retrieve default branch: %v", err.Error())
	}
	cmn.Debug("%s: run: defaultBranch: %v", command, defaultBranch)

	// Define clone path.
	basename := filepath.Base(args[0])
	repo := strings.TrimSuffix(basename, filepath.Ext(basename))
	clonePath := filepath.Join(repo, defaultBranch)

	// Clone the repository.
	err = git.Clone(args[0], defaultBranch, clonePath)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("could not clone repo: %v", err.Error())
	}

	// Write config file to project path.
	err = cmn.WriteConfig(repo, defaultBranch)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error writing config file: %v", err.Error())
	}

	fmt.Printf("Clone complete.\n")

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}
