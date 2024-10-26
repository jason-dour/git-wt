// Package xx implements the xx subcommand for git-wt.
package xx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	} // Cobra command definition for the 'xx' command.
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
	cmn.Debug("%s: %s: begin", command, funcName)

	cmn.Debug("%s: %s: check if most/all used with branches/worktrees", command, funcName)
	if (config.Most || config.All) && (config.Branches || config.Worktrees) {
		if config.Most {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("most and branches/worktrees set; use one or the other")
		}
		if config.All {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("all and branches/worktrees set; use one or the other")
		}
	}

	cmn.Debug("%s: %s: check if most and all both set", command, funcName)
	if config.Most && config.All {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("most and all both set; use one or the other")
	}

	cmn.Debug("%s: %s: check if anything specified to reset", command, funcName)
	if !(config.Branches || config.Worktrees || config.Most || config.All) {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("no flags; must specify what to reset")
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}

// deleteProjectContents deletes all project directory contents except for the configuration file.
func deleteProjectContents() error {
	funcName := "deleteProjectContents"
	cmn.Debug("%s: %s: begin", command, funcName)

	cmn.Debug("%s: %s: reading project directory contents", command, funcName)
	contents, err := os.ReadDir(cmn.Config.ProjectDir)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error reading directory contents: %s", err.Error())
	}
	cmn.Debug("%s: %s: found %d entries in directory", command, funcName, len(contents))

	for _, v := range contents {
		if v.Name() != ".git-wt" {
			if v.Type().IsDir() {
				cmn.Debug("%s: %s: deleting %s", command, funcName, v.Name())
				err := os.RemoveAll(v.Name())
				if err != nil {
					cmn.Debug("%s: %s: error: end", command, funcName)
					return fmt.Errorf("error deleting %s: %s", v.Name(), err.Error())
				}
			} else {
				cmn.Debug("%s: %s: ignoring non-worktree file: %s", command, funcName, v.Name())
			}
		} else {
			cmn.Debug("%s: %s: ignoring config file: %s", command, funcName, v.Name())
		}
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}

// deleteWorktrees deletes all worktrees except for main.
func deleteWorktrees() error {
	funcName := "deleteWorktrees"
	cmn.Debug("%s: %s: begin", command, funcName)

	// Retrieve list of worktrees.
	cmn.Debug("%s: %s: retrieving list of worktrees", command, funcName)
	output, err := git.WorktreeList(true)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error listing worktrees: %s", err.Error())
	}

	// Parse output to build slice of worktrees to delete.
	cmn.Debug("%s: %s: parsing worktree list output", command, funcName)
	worktrees := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		cmn.Debug("%s: %s: output line: %s", command, funcName, line)
		if strings.HasPrefix(line, "worktree ") {
			worktree := filepath.Base(strings.Split(line, " ")[1])
			if filepath.Base(worktree) == cmn.Config.DefaultBranch {
				cmn.Debug("%s: %s: found worktree for default branch: %s; ignoring", command, funcName, worktree)
			} else {
				cmn.Debug("%s: %s: found worktree to delete: %s", command, funcName, worktree)
				worktrees = append(worktrees, worktree)
			}
		}
	}
	cmn.Debug("%s: %s: worktrees to delete: %v", command, funcName, worktrees)

	// Iterate through slice of worktrees and delete them.
	cmn.Debug("%s: %s: deleting worktrees", command, funcName)
	for _, v := range worktrees {
		cmn.Debug("%s: %s: deleting worktree: %s", command, funcName, v)
		_, err := git.WorktreeRemove(&cmn.CfgRm{Force: true}, v)
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error removing worktree: %s", err.Error())
		}
		fmt.Printf("Deleted worktree: %s\n", v)
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}

// deleteBranches deletes all local branches except for main.
func deleteBranches() error {
	funcName := "deleteBranches"
	cmn.Debug("%s: %s: begin", command, funcName)

	branches, err := git.GetBranches()
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error getting branches: %s", err.Error())
	}
	cmn.Debug("%s: %s: branches: %v", command, funcName, branches)

	cmn.Debug("%s: %s: retrieving worktrees", command, funcName)
	worktrees, err := git.WorktreeList(true)
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("error retrieving worktrees: %s", err.Error())
	}

	worktree_branches := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(string(worktrees)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "branch ") {
			worktree_branches[filepath.Base(strings.Split(line, " ")[1])] = struct{}{}
		}
	}
	cmn.Debug("%s: %s: branches checked out in worktrees: %v", command, funcName, worktree_branches)

	for _, v := range branches {
		if v == cmn.Config.DefaultBranch {
			cmn.Debug("%s: %s: found default branch: %s; ignoring", command, funcName, v)
		} else {
			cmn.Debug("%s: %s: deleting branch: %s", command, funcName, v)
			if _, ok := worktree_branches[v]; ok {
				cmn.Debug("%s: %s: error: end", command, funcName)
				return fmt.Errorf("branch checked out in worktree, remove worktree first: %s", v)
			} else {
				err := git.DeleteBranch(v)
				if err != nil {
					cmn.Debug("%s: %s: error: end", command, funcName)
					return fmt.Errorf("error deleting branch: %s", err.Error())
				}
				fmt.Printf("Deleted branch: %s\n", v)
			}
		}
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}

// run is the main function for the 'xx' command.
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

	cmn.Debug("%s: %s: checking if in safe directory", command, funcName)
	if cmn.Config.InitialDir != cmn.Config.ProjectDir {
		// Not in ProjectDir; unsafe.
		cmn.Debug("%s: %s: error: end", command, funcName)
		return fmt.Errorf("not in project dir: change dir to %s", cmn.Config.ProjectDir)
	}
	cmn.Debug("%s: %s: in project directory, proceeding", command, funcName)

	// Check configuration.
	err = checkConfig()
	if err != nil {
		cmn.Debug("%s: %s: error: end", command, funcName)
		return err
	}

	// Delete everything and clone again.
	if config.All {
		cmn.Debug("%s: %s: deleting everything and cloning", command, funcName)

		cmn.Debug("%s: %s: determine remote for cloning after delete", command, funcName)
		remote, err := git.GetRemote(filepath.Join(cmn.Config.InitialDir, cmn.Config.DefaultBranch))
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error determining remote to clone after reset: %s", err.Error())
		}
		cmn.Debug("%s: %s: remote: %s\n", command, funcName, remote)

		cmn.Debug("%s: %s: deleting contents of project directory", command, funcName)
		err = deleteProjectContents()
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error deleting project contents: %s", err.Error())
		}

		cmn.Debug("%s: %s: cloning remote: %s", command, funcName, remote)
		err = git.Clone(remote, cmn.Config.DefaultBranch, cmn.Config.DefaultBranch)
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error cloning remote: %s", err.Error())
		}
	}

	// Delete worktrees.
	if config.Worktrees || config.Most {
		cmn.Debug("%s: %s: deleting worktrees", command, funcName)
		err := deleteWorktrees()
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error deleting worktrees: %s", err.Error())
		}
	}

	// Delete branches.
	if config.Branches || config.Most {
		cmn.Debug("%s: %s: deleting branches", command, funcName)
		err := deleteBranches()
		if err != nil {
			cmn.Debug("%s: %s: error: end", command, funcName)
			return fmt.Errorf("error deleting branches: %s", err.Error())
		}
	}

	cmn.Debug("%s: %s: end", command, funcName)
	return nil
}
