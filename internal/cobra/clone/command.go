package clone

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command = "clone" // Command name.
	Cmd     = &cobra.Command{
		Use:     command + " repo_url",
		Short:   "Clone a repo for a git-wt workflow.",
		Long:    cmn.Basename + " " + command + " - Clone a repo for a git-wt workflow.",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"cl", "clo", "clon"},
		RunE:    run,
	} // Cobra command definition for the 'clone' command.
)

func run(cmd *cobra.Command, args []string) error {
	funcName := "run"
	cmn.Debug("%s: %s: begin\n", command, funcName)
	cmn.Debug("%s: %s: args: %v\n", command, funcName, args)

	fmt.Printf("Cloning %s.\n", args[0])

	// Get default branch from remote repository.
	cmn.Debug("%s: %s: retrieving default branch from remote\n", command, funcName)
	defaultBranch, err := git.GetDefaultBranch(args[0])
	if err != nil {
		cmn.Exit(11, "could not retrieve default branch: %v", err.Error())
	}
	cmn.Debug("%s: run: defaultBranch: %v\n", command, defaultBranch)

	// Define clone path.
	basename := filepath.Base(args[0])
	repo := strings.TrimSuffix(basename, filepath.Ext(basename))
	clonePath := filepath.Join(repo, defaultBranch)

	// Clone the repository.
	err = git.Clone(args[0], defaultBranch, clonePath)
	if err != nil {
		cmn.Exit(12, "could not clone repo: %v", err.Error())
	}

	// Write config file to project path.
	err = cmn.WriteConfig(repo, defaultBranch)
	if err != nil {
		cmn.Exit(13, "error writing config file: %v", err.Error())
	}

	fmt.Printf("Clone complete.\n")

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}