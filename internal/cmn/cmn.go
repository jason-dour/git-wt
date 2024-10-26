// Package cmn implements common variables and utility functions for git-wt,
// providing debug and configuration.
package cmn

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type (
	Cfg struct {
		DebugFlag     bool   // Whether debug output is enabled.
		DefaultBranch string // Default branch/worktree of the project.
		InitialDir    string // Initial working directory of the program.
		ProjectDir    string // Path of the project directory.
	} // Configuration for the program.

	CfgMk struct {
		Branch      string
		BranchReset string
		CheckoutNo  bool
		Force       bool
		Track       bool
		Quiet       bool
	} // Configuration for 'mk' command.

	CfgMv struct {
		Force bool
	} // Configuration for 'mv' command.

	CfgRm struct {
		Force bool
	} // Configuration for 'rm' command.

	CfgXx struct {
		Branches  bool // Whether to reset branches.
		Worktrees bool // Whether to reset worktrees.
		Most      bool // Whether to reset both branches and worktrees.
		All       bool // Whether to wipe everyting and clone again.
	} // Configuration for 'xx' command.
)

var (
	Basename string          // Base name of the program; injected during compile.
	Version  string          // Version of the program; injected during compile.
	Config   *Cfg   = &Cfg{} // Global configuration for the program.
)

// Debug writes debug output to Stderr if DebugFlag is true.
func Debug(format string, args ...interface{}) {
	if Config.DebugFlag {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("debug: "+format, args...))
	}
}

// WriteConfig writes program's config file to cloned repo's project path.
func WriteConfig(path string, branch string) error {
	funcName := "cmn.WriteConfig"
	Debug("%s: begin", funcName)

	// Set the configuration file name.
	filename := filepath.Join(path, "."+Basename)
	Debug("%s: config filename: %s", funcName, filename)

	// Open the configuration file to write the config.
	cfgFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		Debug("%s: error: end", funcName)
		return fmt.Errorf("could not create config file: %v", err.Error())
	}
	defer cfgFile.Close()
	Debug("%s: config file opened for writing", funcName)

	// Write the configuration to the file.
	cfg := fmt.Sprintf("default: %s", branch)
	fmt.Fprintf(cfgFile, "%s", cfg)
	Debug("%s: config written to file> %s", funcName, cfg)

	Debug("%s: end", funcName)
	return nil
}

// FindConfig finds the program's config file.
func findConfig(cwd string) (string, error) {
	funcName := "cmn.findConfig"
	Debug("%s: begin", funcName)

	cfgFile, err := walkUpTree(cwd)
	if err != nil {
		Debug("%s: error: end", funcName)
		return "", fmt.Errorf("no project config file found")
	}

	Debug("%s: end", funcName)
	return cfgFile, nil
}

// walkUpTree walks up the filesystem tree until the config file is found.
func walkUpTree(path string) (string, error) {
	funcName := "cmn.walkUpTree"
	Debug("%s: begin", funcName)

	target := filepath.Clean(path)
	for {
		filename := filepath.Join(target, "."+Basename)
		Debug("%s: checking %s", funcName, filename)

		_, err := os.Stat(filename)
		if err != nil {
			target, _ = filepath.Split(target)
			target = filepath.Clean(target)
			if target == "/" {
				Debug("%s: error: end", funcName)
				return "", fmt.Errorf("no project config file")
			}
			continue
		} else {
			Config.ProjectDir = target
			Debug("%s: project dir: %s", funcName, Config.ProjectDir)
			Debug("%s: end", funcName)
			return filename, nil
		}
	}
}

// InitConfig locates and reads the configuration file for use.
func InitConfig() error {
	funcName := "cmn.InitConfig"
	Debug("%s: begin", funcName)

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		Debug("%s: error: end", funcName)
		return fmt.Errorf("%s: error locating working directory: %s", funcName, err.Error())
	}
	Config.InitialDir = cwd
	Debug("%s: set initial dir: %s", funcName, Config.InitialDir)

	// Locate the config file.
	cfgFile, err := findConfig(Config.InitialDir)
	if err != nil {
		Debug("%s: error: end\n", funcName)
		return fmt.Errorf("%s: error locating config: %s", funcName, err.Error())
	}
	Debug("%s: config file: %s", funcName, cfgFile)

	// Read the config file.
	file, err := os.ReadFile(cfgFile)
	if err != nil {
		Debug("%s: error: end", funcName)
		return fmt.Errorf("%s: error reading config file: %s", funcName, err.Error())
	}
	Debug("%s: config file length: %d", funcName, len(file))
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	for scanner.Scan() {
		line := scanner.Text()
		Debug("%s: config file line: %s", funcName, line)
		if strings.HasPrefix(line, "default: ") {
			Config.DefaultBranch = strings.Split(line, " ")[1]
			break
		} else {
			Debug("%s: invalid config line; skipping: %s", funcName, line)
		}
	}
	Debug("%s: default branch: %s", funcName, Config.DefaultBranch)

	Debug("%s: end", funcName)
	return nil
}
