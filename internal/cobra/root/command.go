// Package root implements the root command for git-wt.
package root

import (
	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/cobra/cl"
	"github.com/jason-dour/git-wt/internal/cobra/ls"
	"github.com/jason-dour/git-wt/internal/cobra/mk"
	"github.com/jason-dour/git-wt/internal/cobra/mv"
	"github.com/jason-dour/git-wt/internal/cobra/rm"
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
	// Command flags.
	Cmd.PersistentFlags().BoolVarP(&cmn.Config.DebugFlag, "debug", "d", false, "enable debug mode")

	// Sub-Commands
	Cmd.AddCommand(cl.Cmd)
	Cmd.AddCommand(ls.Cmd)
	Cmd.AddCommand(mk.Cmd)
	Cmd.AddCommand(mv.Cmd)
	Cmd.AddCommand(rm.Cmd)
}
