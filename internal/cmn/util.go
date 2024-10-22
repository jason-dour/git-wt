package cmn

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Debug writes debug output to Stderr if DebugFlag is true.
func Debug(format string, args ...interface{}) {
	if DebugFlag {
		fmt.Fprintf(os.Stderr, "debug: "+format, args...)
	}
}

// Exit writes the provided error details to Stderr then exits with the provided value.
func Exit(retval int, format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("error: "+format, args...))
	os.Exit(retval)
}

// WriteConfig writes program's config file to cloned repo's project path.
func WriteConfig(path string, branch string) error {
	funcName := "WriteConfig"
	Debug("%s: begin\n", funcName)

	// Set the configuration file name.
	filename := filepath.Join(path, "."+Basename)
	Debug("%s: config filename: %s\n", funcName, filename)

	// Open the configuration file to write the config.
	cfgFile, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_SYNC|os.O_RDWR, 0644)
	if err != nil {
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
			ProjectDir = target
			Debug("%s: project dir: %s", funcName, ProjectDir)
			Debug("%s: end\n", funcName)
			return filename, nil
		}
	}
}

func InitConfig() {
	funcName := "cmn.initConfig"
	Debug("%s: begin\n", funcName)

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		Exit(127, "error locating working directory: %s", err.Error())
	}
	InitialDir = cwd
	Debug("%s: set initial dir: %s\n", funcName, InitialDir)

	// Locate the config file.
	cfgFile, _ := findConfig(InitialDir)
	// if err != nil {
	// 	Exit(127, "error locating config: %s", err.Error())
	// }
	Debug("%s: config file: %s\n", funcName, cfgFile)

	// Read the config file.
	file, _ := os.ReadFile(cfgFile)
	// if err != nil {
	// 	Exit(127, "error reading config file: %s", err.Error())
	// }
	Debug("%s: config file contents > %s\n", funcName, string(file))
	if file != nil {
		if slices.Equal(file[0:9], []byte("default: ")) {
			DefaultBranch = string(file[9:])
		} else {
			// Exit(127, "error loading config file")
		}
	}
	Debug("%s: default branch: %s\n", funcName, DefaultBranch)

	Debug("%s: end\n", funcName)
}
