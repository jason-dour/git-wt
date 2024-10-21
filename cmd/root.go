// cmd/root
//
// Root command for git-wt.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version  string // The version of the program; injected at compile time.
	basename string // The basename of the program; injected at compile time.

	rootCmd = &cobra.Command{
		Use:   basename,
		Short: "A git worktree helper.",
		Long: `Git-wt is a git helper for managing worktrees as part of
a worktree-based workflow.`,
		Run:     func(cmd *cobra.Command, args []string) { run() },
		Version: version,
	}
)

func init() {
	// Initialize Cobra.
	// cobra.OnInitialize(initConfig)

	// Define the command line flags.
	// rootCmd.PersistentFlags()...
}

func run() {
	fmt.Printf("nothing here just yet.\n")
}

func Execute() error {
	return rootCmd.Execute()
}
