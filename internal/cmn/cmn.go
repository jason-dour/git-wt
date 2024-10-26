// Package cmn implements common variables and utility functions for git-wt,
// providing debug and configuration.
package cmn

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

	CfgRm struct {
		Force bool
	} // Configuration for 'rm' command.

	CfgMv struct {
		Force bool
	} // Configuration for 'mv' command.
)

var (
	Basename string          // Base name of the program; injected during compile.
	Version  string          // Version of the program; injected during compile.
	Config   *Cfg   = &Cfg{} // Global configuration for the program.
)

// Debug writes debug output to Stderr if DebugFlag is true.
func Debug(format string, args ...interface{}) {
	if Config.DebugFlag {
		fmt.Fprintf(os.Stderr, "debug: "+format, args...)
	}
}

// WriteConfig writes program's config file to cloned repo's project path.
func WriteConfig(path string, branch string) error {
	funcName := "cmn.WriteConfig"
	Debug("%s: begin\n", funcName)

	// Set the configuration file name.
	filename := filepath.Join(path, "."+Basename)
	Debug("%s: config filename: %s\n", funcName, filename)

	// Open the configuration file to write the config.
	cfgFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		Debug("%s: end\n", funcName)
		return fmt.Errorf("could not create config file: %v", err.Error())
	}
	defer cfgFile.Close()
	Debug("%s: config file opened for writing\n", funcName)

	// Write the configuration to the file.
	cfg := fmt.Sprintf("default: %s", branch)
	fmt.Fprintf(cfgFile, "%s", cfg)
	Debug("%s: config written to file> %s\n", funcName, cfg)

	Debug("%s: end\n", funcName)
	return nil
}

// FindConfig finds the program's config file.
func findConfig(cwd string) (string, error) {
	funcName := "cmn.findConfig"
	Debug("%s: begin\n", funcName)

	cfgFile, err := walkUpTree(cwd)
	if err != nil {
		Debug("%s: end\n", funcName)
		return "", fmt.Errorf("no project config file found")
	}

	Debug("%s: end\n", funcName)
	return cfgFile, nil
}

// walkUpTree walks up the filesystem tree until the config file is found.
func walkUpTree(path string) (string, error) {
	funcName := "cmn.walkUpTree"
	Debug("%s: begin\n", funcName)

	target := filepath.Clean(path)
	for {
		// Debug("%s: looking in %s\n", funcName, target)
		filename := filepath.Join(target, "."+Basename)
		Debug("%s: checking %s\n", funcName, filename)

		_, err := os.Stat(filename)
		if err != nil {
			target, _ = filepath.Split(target)
			target = filepath.Clean(target)
			if target == "/" {
				Debug("%s: end\n", funcName)
				return "", fmt.Errorf("no project config file")
			}
			continue
		} else {
			Config.ProjectDir = target
			Debug("%s: project dir: %s\n", funcName, Config.ProjectDir)
			Debug("%s: end\n", funcName)
			return filename, nil
		}
	}
}

// InitConfig locates and reads the configuration file for use.
func InitConfig() error {
	funcName := "cmn.InitConfig"
	Debug("%s: begin\n", funcName)

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		Debug("%s: end\n", funcName)
		return fmt.Errorf("%s: error locating working directory: %s", funcName, err.Error())
	}
	Config.InitialDir = cwd
	Debug("%s: set initial dir: %s\n", funcName, Config.InitialDir)

	// Locate the config file.
	cfgFile, err := findConfig(Config.InitialDir)
	if err != nil {
		Debug("%s: end\n", funcName)
		return fmt.Errorf("%s: error locating config: %s", funcName, err.Error())
	}
	Debug("%s: config file: %s\n", funcName, cfgFile)

	// Read the config file.
	file, err := os.ReadFile(cfgFile)
	if err != nil {
		Debug("%s: end\n", funcName)
		return fmt.Errorf("%s: error reading config file: %s", funcName, err.Error())
	}
	Debug("%s: config file contents > %s\n", funcName, string(file))
	if file != nil {
		if slices.Equal(file[0:9], []byte("default: ")) {
			Config.DefaultBranch = string(file[9:])
		} else {
			Debug("%s: end\n", funcName)
			return fmt.Errorf(":%s: error parsing config file: invalid format", funcName)
		}
	}
	Debug("%s: default branch: %s\n", funcName, Config.DefaultBranch)

	Debug("%s: end\n", funcName)
	return nil
}
