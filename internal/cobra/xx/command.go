// Package xx implements the xx subcommand for git-wt.
package xx

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jason-dour/git-wt/internal/cmn"
	"github.com/jason-dour/git-wt/internal/git"
	"github.com/spf13/cobra"
)

var (
	command            = "xx"         // Command name.
	config  *cmn.CfgXx = &cmn.CfgXx{} // Configuration for the command.
	Cmd                = &cobra.Command{
		Use:     command,
		Short:   "Reset project.",
		Long:    cmn.Basename + " " + command + " - Reset project.",
		Args:    cobra.NoArgs,
		Aliases: []string{"list"},
		RunE:    run,
		Hidden:  true,
	} // Cobra command definition for the 'list' command.
)

// init performs initialization for the 'mk' command.
func init() {
	Cmd.PersistentFlags().BoolVarP(&config.Branches, "branches", "b", false, "delete local branches")
	Cmd.PersistentFlags().BoolVarP(&config.Worktrees, "worktrees", "w", false, "delete all worktrees except main")
	Cmd.PersistentFlags().BoolVarP(&config.Most, "most", "m", false, "delete local branches and all worktrees except main")
	Cmd.PersistentFlags().BoolVarP(&config.All, "all", "a", false, "delete everything and clone again")
}

// checkConfig scans config for proper use of flags.
func checkConfig() error {
	funcName := "checkConfig"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	cmn.Debug("%s: %s: check if most/all used with branches/worktrees\n", command, funcName)
	if (config.Most || config.All) && (config.Branches || config.Worktrees) {
		if config.Most {
			return fmt.Errorf("most and branches/worktrees set; use one or the other")
		}
		if config.All {
			return fmt.Errorf("all and branches/worktrees set; use one or the other")
		}
	}

	cmn.Debug("%s: %s: check if most and all both set\n", command, funcName)
	if config.Most && config.All {
		return fmt.Errorf("most and all both set; use one or the other")
	}

	cmn.Debug("%s: %s: check if anything specified to reset\n", command, funcName)
	if !(config.Branches || config.Worktrees || config.Most || config.All) {
		return fmt.Errorf("no flags; must specify what to reset")
	}

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}

// deleteProjectContents deletes all project directory contents except for the configuration file.
func deleteProjectContents() error {
	funcName := "deleteProjectContents"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	cmn.Debug("%s: %s: reading project directory contents\n", command, funcName)
	contents, err := os.ReadDir(cmn.Config.ProjectDir)
	if err != nil {
		return fmt.Errorf("error reading directory contents: %s", err.Error())
	}
	cmn.Debug("%s: %s: found %d entries in directory\n", command, funcName, len(contents))

	for _, v := range contents {
		if v.Name() != ".git-wt" {
			if v.Type().IsDir() {
				cmn.Debug("%s: %s: deleting %s\n", command, funcName, v.Name())
				err := os.RemoveAll(v.Name())
				if err != nil {
					return fmt.Errorf("error deleting %s: %s", v.Name(), err.Error())
				}
			} else {
				cmn.Debug("%s: %s: ignoring non-worktree file: %s\n", command, funcName, v.Name())
			}
		} else {
			cmn.Debug("%s: %s: ignoring config file: %s\n", command, funcName, v.Name())
		}
	}

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}

// deleteWorktrees deletes all worktrees except for main.
func deleteWorktrees() error {
	funcName := "deleteWorktrees"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}

// deleteBranches deletes all local branches except for main.
func deleteBranches() error {
	funcName := "deleteBranches"
	cmn.Debug("%s: %s: begin\n", command, funcName)

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}

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

	cmn.Debug("%s: %s: checking if in safe directory\n", command, funcName)
	if cmn.Config.InitialDir != cmn.Config.ProjectDir {
		// Not in ProjectDir; unsafe.
		return fmt.Errorf("not in project dir: change dir to %s", cmn.Config.ProjectDir)
	}
	cmn.Debug("%s: %s: in project directory, proceeding\n", command, funcName)

	// Check configuration.
	err = checkConfig()
	if err != nil {
		return err
	}

	// Delete everything and clone again.
	if config.All {
		cmn.Debug("%s: %s: deleting everything and cloning\n", command, funcName)

		cmn.Debug("%s: %s: determine remote for cloning after delete\n", command, funcName)
		remote, err := git.GetRemote(filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch))
		if err != nil {
			return fmt.Errorf("error determining remote to clone after reset: %s", err.Error())
		}
		cmn.Debug("%s: %s: remote: %s\n", command, funcName, remote)

		cmn.Debug("%s: %s: deleting contents of project directory\n", command, funcName)
		err = deleteProjectContents()
		if err != nil {
			return fmt.Errorf("error deleting project contents: %s", err.Error())
		}

		cmn.Debug("%s: %s: cloning remote: %s\n", command, funcName, remote)
		err = git.Clone(remote, cmn.Config.DefaultBranch, cmn.Config.DefaultBranch)
		if err != nil {
			return fmt.Errorf("error cloning remote: %s", err.Error())
		}
	}

	// Delete worktrees.
	if config.Worktrees || config.Most {
		cmn.Debug("%s: %s: deleting worktrees\n", command, funcName)
		err := deleteWorktrees()
		if err != nil {
			return fmt.Errorf("error deleting worktrees: %s", err.Error())
		}
	}

	// Delete branches.
	if config.Branches || config.Most {
		cmn.Debug("%s: %s: deleting branches\n", command, funcName)
		err := deleteBranches()
		if err != nil {
			return fmt.Errorf("error deleting branches: %s", err.Error())
		}
	}

	cmn.Debug("%s: %s: end\n", command, funcName)
	return nil
}
