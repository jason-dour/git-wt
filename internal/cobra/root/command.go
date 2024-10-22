// cmd/root
//
// Root command for git-wt.

package root

import (
	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/cobra/clone"
	"github.com/jason-dour/git-wt/internal/cobra/list"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     cmn.Basename,
		Long:    cmn.Basename + " - A git helper for managing worktrees as part of a workflow.",
		Args:    cobra.NoArgs,
		Version: cmn.Version,
	} // Cobra root command definition.
)

func init() {
	// Load configuration.
	cobra.OnInitialize(cmn.InitConfig)

	// Command flags.
	Cmd.PersistentFlags().BoolVarP(&cmn.DebugFlag, "debug", "d", false, "enable debug mode")

	// Sub-Commands
	Cmd.AddCommand(clone.Cmd)
	Cmd.AddCommand(list.Cmd)
}
