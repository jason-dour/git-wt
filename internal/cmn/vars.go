package cmn

var (
	Basename string // Base name of the program; injected during compile.
	Version  string // Version of the program; injected during compile.

	DebugFlag bool // Whether debug output is enabled.

	InitialDir    string // Initial working directory of the program.
	ProjectDir    string // Path of the project directory.
	DefaultBranch string // Default branch/worktree of the project.
)
