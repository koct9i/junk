package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// https://pkg.go.dev/cmd/go
// https://cs.opensource.google/go/go/+/master:src/cmd/go/alldocs.go

// buildFlags are the flags shared by build, clean, get, install, list, run, and test.
func buildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "C", Usage: "Change to 'dir' before running the command."},
		&cli.BoolFlag{Name: "a", Usage: "Force rebuilding of packages that are already up-to-date."},
		&cli.BoolFlag{Name: "n", Usage: "Print the commands but do not run them."},
		&cli.IntFlag{Name: "p", Value: 0, Usage: "The number of programs that can be run in parallel."},
		&cli.BoolFlag{Name: "race", Usage: "Enable data race detection."},
		&cli.BoolFlag{Name: "msan", Usage: "Enable interoperation with memory sanitizer."},
		&cli.BoolFlag{Name: "asan", Usage: "Enable interoperation with address sanitizer."},
		&cli.BoolFlag{Name: "cover", Usage: "Enable code coverage instrumentation."},
		&cli.StringFlag{Name: "covermode", Usage: "Coverage mode: set, count, atomic."},
		&cli.StringSliceFlag{Name: "coverpkg", Usage: "Comma-separated patterns of packages for which to apply coverage in main builds."},
		&cli.BoolFlag{Name: "v", Usage: "Print the names of packages as they are compiled."},
		&cli.BoolFlag{Name: "work", Usage: "Print the name of the temporary work directory and do not delete it."},
		&cli.BoolFlag{Name: "x", Usage: "Print the commands."},
		&cli.StringSliceFlag{Name: "asmflags", Usage: "Arguments to pass on each 'go tool asm' invocation. Supports optional [pattern=] prefix."},
		&cli.StringFlag{Name: "buildmode", Usage: "Build mode to use. See 'go help buildmode'."},
		&cli.StringFlag{Name: "buildvcs", Usage: "Whether to stamp binaries with VCS info: true, false, or auto."},
		&cli.StringFlag{Name: "compiler", Usage: "Compiler to use: gc or gccgo."},
		&cli.StringSliceFlag{Name: "gccgoflags", Usage: "Arguments to pass on each gccgo compiler/linker invocation."},
		&cli.StringSliceFlag{Name: "gcflags", Usage: "Arguments to pass on each 'go tool compile' invocation. Supports optional [pattern=] prefix."},
		&cli.StringFlag{Name: "installsuffix", Usage: "Suffix to use in the package installation directory."},
		&cli.BoolFlag{Name: "json", Usage: "Emit build output in JSON for automated processing."},
		&cli.StringSliceFlag{Name: "ldflags", Usage: "Arguments to pass on each 'go tool link' invocation."},
		&cli.BoolFlag{Name: "linkshared", Usage: "Build code that will be linked against shared libraries (with -buildmode=shared)."},
		&cli.StringFlag{Name: "mod", Usage: "Module download mode: readonly, vendor, or mod."},
		&cli.BoolFlag{Name: "modcacherw", Usage: "Leave newly-created module cache directories read-write."},
		&cli.StringFlag{Name: "modfile", Usage: "Read (and possibly write) an alternate go.mod file."},
		&cli.StringSliceFlag{Name: "pkgdir", Usage: "Package directories (for gccgo)."},
		&cli.StringSliceFlag{Name: "tags", Usage: "Build tags to consider satisfied during the build."},
		&cli.StringFlag{Name: "toolexec", Usage: "Program to use to invoke toolchain programs (e.g., compile, link)."},
		&cli.StringFlag{Name: "trimpath", Usage: "Remove all file system paths from the resulting executable."},
		&cli.StringFlag{Name: "toolchain", Usage: "Select the Go toolchain to use."},
	}
}

// listFlags are the list-specific flags (in addition to build flags).
func listFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{Name: "deps", Usage: "List dependencies of each package."},
		&cli.StringFlag{Name: "f", Usage: "Print using a custom format."},
		&cli.BoolFlag{Name: "find", Usage: "Show a suggestion of the package pattern, avoiding loading of packages."},
		&cli.BoolFlag{Name: "json", Usage: "Print JSON instead of the default text format."},
		&cli.BoolFlag{Name: "m", Usage: "List modules instead of packages."},
		&cli.BoolFlag{Name: "test", Usage: "Include in the listing test packages imported by each package."},
		&cli.BoolFlag{Name: "compiled", Usage: "Include the names of compiled Go files."},
		&cli.BoolFlag{Name: "cgo", Usage: "Include the names of Cgo files."},
		&cli.BoolFlag{Name: "e", Usage: "Show errors associated with packages."},
		&cli.BoolFlag{Name: "export", Usage: "Include the file names for export data."},
		&cli.StringFlag{Name: "reuse", Usage: "Read previously saved JSON from a 'go list -m -json' run (for -m)."},
		&cli.BoolFlag{Name: "u", Usage: "When -m is set, also show available minor and patch upgrades (with -u/-versions)."},
		&cli.BoolFlag{Name: "retracted", Usage: "When -m is set, include retracted versions."},
		&cli.BoolFlag{Name: "versions", Usage: "When -m is set, show available versions."},
	}
}

// testFlags are the test-execution flags recognized by 'go test'.
func testFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "bench", Usage: "Run only those benchmarks matching the regexp."},
		&cli.StringFlag{Name: "benchtime", Usage: "Run enough iterations of each benchmark to take the specified time (e.g., 1s, 2x)."},
		&cli.BoolFlag{Name: "benchmem", Usage: "Print memory allocation statistics for benchmarks."},
		&cli.IntFlag{Name: "blockprofilerate", Usage: "Control the fraction of goroutine blocking events reported in the blocking profile."},
		&cli.StringFlag{Name: "blockprofile", Usage: "Write a goroutine blocking profile to the named file after execution."},
		&cli.StringFlag{Name: "coverprofile", Usage: "Write a coverage profile to the named file after execution."},
		&cli.StringFlag{Name: "cpu", Usage: "Comma-separated list of GOMAXPROCS values for which the tests should be executed."},
		&cli.IntFlag{Name: "count", Value: 1, Usage: "Run each test and benchmark n times (0 means infinite)."},
		&cli.BoolFlag{Name: "failfast", Usage: "Do not start new tests after the first test failure."},
		&cli.BoolFlag{Name: "fullpath", Usage: "Show full pathnames in test output."},
		&cli.StringFlag{Name: "list", Usage: "List tests, benchmarks, or examples matching the regexp and exit."},
		&cli.StringFlag{Name: "memprofile", Usage: "Write a memory profile to the named file after execution."},
		&cli.IntFlag{Name: "memprofilerate", Usage: "Enable more precise memory profiles (see runtime.MemProfileRate)."},
		&cli.StringFlag{Name: "mutexprofile", Usage: "Write a mutex contention profile to the named file after execution."},
		&cli.IntFlag{Name: "mutexprofilefraction", Usage: "Sample fraction of mutex contention events."},
		&cli.StringFlag{Name: "outputdir", Usage: "Write profiles to the specified directory, leaving the test binary."},
		&cli.IntFlag{Name: "parallel", Usage: "Maximum number of tests to run in parallel."},
		&cli.StringFlag{Name: "run", Usage: "Run only those tests and examples matching the regexp."},
		&cli.BoolFlag{Name: "short", Usage: "Tell long-running tests to shorten their run time."},
		&cli.DurationFlag{Name: "timeout", Usage: "If a test runs longer than t, panic. The default is 10m."},
		&cli.BoolFlag{Name: "trace", Usage: "Write an execution trace to the named file after execution (when used with value)."},
		&cli.BoolFlag{Name: "v", Usage: "Verbose output: log all tests as they are run."},
	}
}

// ---------- Commands (no-op handlers) ----------

func stubAction(ctx context.Context, command *cli.Command) error {
	fmt.Printf("%s: not implemented\n", command.Name)
	return nil
}

func cmdBug() *cli.Command {
	return &cli.Command{
		Name:   "bug",
		Usage:  "Start a bug report",
		Action: stubAction,
	}
}

func cmdBuild() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "Compile packages and dependencies",
		Flags:  buildFlags(),
		Action: stubAction,
	}
}

func cmdClean() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags,
		&cli.BoolFlag{Name: "cache", Usage: "Remove all cached build and test results."},
		&cli.BoolFlag{Name: "i", Usage: "Remove the installed packages for the named targets."},
		&cli.BoolFlag{Name: "modcache", Usage: "Remove the entire module download cache."},
		&cli.BoolFlag{Name: "r", Usage: "Remove obj and installed files recursively for the arguments and their dependencies."},
		&cli.BoolFlag{Name: "testcache", Usage: "Expire all test results in the cache."},
	)
	return &cli.Command{
		Name:   "clean",
		Usage:  "Remove object files and cached files",
		Flags:  flags,
		Action: stubAction,
	}
}

func cmdDoc() *cli.Command {
	return &cli.Command{
		Name:  "doc",
		Usage: "Show documentation for package or symbol",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "all", Usage: "Show documentation for all symbols, as per 'go doc -all'."},
			&cli.StringFlag{Name: "c", Usage: "Respect case when matching symbols (-c)."}, // urfave/cli requires a value; keep present for parity
			&cli.StringFlag{Name: "cmd", Usage: "Treat a cmd path as a package (show package docs)."},
			&cli.BoolFlag{Name: "u", Usage: "Show documentation for unexported symbols."},
			&cli.StringFlag{Name: "src", Usage: "Show the source code for the symbol."},
		},
		Action: stubAction,
	}
}

func cmdEnv() *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "Print Go environment information",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json"},
			&cli.BoolFlag{Name: "changed"},
			&cli.BoolFlag{Name: "u"},
			&cli.StringSliceFlag{Name: "w", Usage: "Set NAME=VALUE (may be repeated)."},
		},
		Action: stubAction,
	}
}

func cmdFix() *cli.Command {
	return &cli.Command{Name: "fix", Usage: "Update packages to use new APIs", Flags: []cli.Flag{
		&cli.StringFlag{Name: "fix", Usage: "Comma-separated list of fixes to run (default all)."},
	}, Action: stubAction}
}

func cmdFmt() *cli.Command {
	return &cli.Command{Name: "fmt", Usage: "Gofmt (reformat) package sources", Flags: []cli.Flag{
		&cli.BoolFlag{Name: "n"},
		&cli.BoolFlag{Name: "x"},
	}, Action: stubAction}
}

func cmdGenerate() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags,
		&cli.StringFlag{Name: "run", Usage: "Run only generators matching the regular expression"},
		&cli.StringFlag{Name: "tags", Usage: "Build tags (also respected by generate)."},
		&cli.BoolFlag{Name: "n", Usage: "Print commands but do not run them."},
		&cli.BoolFlag{Name: "v", Usage: "Verbose output."},
	)
	return &cli.Command{Name: "generate", Usage: "Generate Go files by processing source", Flags: flags, Action: stubAction}
}

func cmdGet() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags,
		&cli.BoolFlag{Name: "d", Usage: "Only download; do not build or install."},
		&cli.BoolFlag{Name: "t", Usage: "Also download test dependencies."},
		&cli.BoolFlag{Name: "u", Usage: "Update modules providing dependencies to newer minor or patch releases; use -u=patch for patch only."},
		&cli.StringFlag{Name: "u=patch", Usage: "Shorthand for patch-level updates only (mirror of go help get)."},
	)
	return &cli.Command{Name: "get", Usage: "Add dependencies to current module and install them", Flags: flags, Action: stubAction}
}

func cmdInstall() *cli.Command {
	return &cli.Command{
		Name:   "install",
		Usage:  "Compile and install packages and dependencies",
		Flags:  buildFlags(),
		Action: stubAction,
	}
}

func cmdList() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags, listFlags()...)
	return &cli.Command{
		Name:   "list",
		Usage:  "List packages or modules",
		Flags:  flags,
		Action: stubAction,
	}
}

func cmdRun() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags,
		&cli.StringFlag{Name: "exec", Usage: "Run the generated binary under xprog (like 'time')."},
	)
	return &cli.Command{
		Name:   "run",
		Usage:  "Compile and run a Go program",
		Flags:  flags,
		Action: stubAction,
	}
}

func cmdTelemetry() *cli.Command {
	return &cli.Command{
		Name:  "telemetry",
		Usage: "Manage telemetry data and settings",
		Commands: []*cli.Command{
			{Name: "on", Usage: "Enable both local collection and uploading.", Action: stubAction},
			{Name: "off", Usage: "Disable both collection and uploading.", Action: stubAction},
			{Name: "local", Usage: "Disable uploading, keep local collection.", Action: stubAction},
		},
		Action: stubAction,
	}
}

func cmdTest() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	flags = append(flags, testFlags()...)
	return &cli.Command{
		Name:   "test",
		Usage:  "Test packages",
		Flags:  flags,
		Action: stubAction,
	}
}

func cmdTool() *cli.Command {
	return &cli.Command{
		Name:   "tool",
		Usage:  "Run specified go tool",
		Action: stubAction,
	}
}

func cmdVersion() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Print Go version",
		Action: stubAction,
	}
}

func cmdVet() *cli.Command {
	flags := append([]cli.Flag{}, buildFlags()...)
	return &cli.Command{
		Name:   "vet",
		Usage:  "Report likely mistakes in packages",
		Flags:  flags,
		Action: stubAction,
	}
}

func cmdMod() *cli.Command {
	return &cli.Command{
		Name:  "mod",
		Usage: "Module maintenance",
		Commands: []*cli.Command{
			{
				Name:   "download",
				Usage:  "Download modules to local cache",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "x"}, &cli.BoolFlag{Name: "json"}, &cli.StringFlag{Name: "reuse"}},
				Action: stubAction,
			},
			{
				Name:  "edit",
				Usage: "Edit go.mod from tools or scripts",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "fmt"},
					&cli.StringSliceFlag{Name: "require", Usage: "Add a requirement (path@version)."},
					&cli.StringSliceFlag{Name: "droprequire", Usage: "Drop a requirement (path)."},
					&cli.StringSliceFlag{Name: "exclude", Usage: "Add an exclude directive (path@version)."},
					&cli.StringSliceFlag{Name: "dropexclude", Usage: "Drop an exclude directive (path@version)."},
					&cli.StringSliceFlag{Name: "replace", Usage: "Add a replace directive old[@v]=new[@v]."},
					&cli.StringSliceFlag{Name: "dropreplace", Usage: "Drop a replace directive old[@v]."},
					&cli.StringFlag{Name: "retract", Usage: "Add a retract directive."},
					&cli.StringFlag{Name: "dropretract", Usage: "Drop a retract directive."},
					&cli.StringFlag{Name: "go", Usage: "Set the expected Go language version."},
					&cli.StringFlag{Name: "toolchain", Usage: "Set the toolchain to use."},
					&cli.BoolFlag{Name: "print"},
					&cli.BoolFlag{Name: "json"},
				},
				Action: stubAction,
			},
			{
				Name:   "graph",
				Usage:  "Print module requirement graph",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "x"}},
				Action: stubAction,
			},
			{
				Name:   "init",
				Usage:  "Initialize new module in current directory",
				Flags:  []cli.Flag{&cli.StringFlag{Name: "module-path", Usage: "Optional module path"}},
				Action: stubAction,
			},
			{
				Name:  "tidy",
				Usage: "Add missing and remove unused modules",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "compat", Usage: "(via -compat=) use compatibility with previous Go version (string value in go tool)."},
					&cli.BoolFlag{Name: "e", Usage: "Report errors but proceed (best-effort)."},
					&cli.BoolFlag{Name: "v"},
					&cli.BoolFlag{Name: "x"},
				},
				Action: stubAction,
			},
			{
				Name:   "vendor",
				Usage:  "Make vendored copy of dependencies",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "e"}, &cli.BoolFlag{Name: "v"}, &cli.StringFlag{Name: "o", Usage: "Output directory"}},
				Action: stubAction,
			},
			{
				Name:   "verify",
				Usage:  "Verify dependencies have expected content",
				Action: stubAction,
			},
			{
				Name:   "why",
				Usage:  "Explain why packages or modules are needed",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "m"}, &cli.BoolFlag{Name: "vendor"}},
				Action: stubAction,
			},
		},
	}
}

func cmdWork() *cli.Command {
	return &cli.Command{
		Name:  "work",
		Usage: "Workspace maintenance",
		Commands: []*cli.Command{
			{
				Name:  "edit",
				Usage: "Edit go.work from tools or scripts",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "fmt"},
					&cli.StringSliceFlag{Name: "godebug", Usage: "Add godebug key=value (may repeat)."},
					&cli.StringSliceFlag{Name: "dropgodebug", Usage: "Drop godebug key."},
					&cli.StringSliceFlag{Name: "use", Usage: "Add use=path directive (may repeat)."},
					&cli.StringSliceFlag{Name: "dropuse", Usage: "Drop use=path directive (may repeat)."},
					&cli.StringSliceFlag{Name: "replace", Usage: "Add replace old[@v]=new[@v]."},
					&cli.StringSliceFlag{Name: "dropreplace", Usage: "Drop replace old[@v]."},
					&cli.StringFlag{Name: "go", Usage: "Set expected Go language version."},
					&cli.StringFlag{Name: "toolchain", Usage: "Set toolchain name."},
					&cli.BoolFlag{Name: "print"},
					&cli.BoolFlag{Name: "json"},
				},
				Action: stubAction,
			},
			{
				Name:   "init",
				Usage:  "Initialize workspace file",
				Action: stubAction,
			},
			{
				Name:   "sync",
				Usage:  "Sync workspace build list to modules",
				Action: stubAction,
			},
			{
				Name:   "use",
				Usage:  "Add modules to workspace file",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "r", Usage: "Search directories recursively."}},
				Action: stubAction,
			},
			{
				Name:   "vendor",
				Usage:  "Make vendored copy of dependencies",
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "e"}, &cli.BoolFlag{Name: "v"}, &cli.StringFlag{Name: "o", Usage: "Output directory"}},
				Action: stubAction,
			},
		},
	}
}

func cmdGo() *cli.Command {
	return &cli.Command{
		Name:      "go",
		Usage:     "shell completion for go command",
		UsageText: "go <command> [arguments...]",
		Description: `Go is a tool for managing Go source code.

Documentation: https://pkg.go.dev/cmd/go`,
		HideHelp: true,
		EnableShellCompletion: true,
		ConfigureShellCompletionCommand: func(completion *cli.Command) {
			completion.Hidden = false
		},
		Commands: []*cli.Command{
			cmdBug(),
			cmdBuild(),
			cmdClean(),
			cmdDoc(),
			cmdEnv(),
			cmdFix(),
			cmdFmt(),
			cmdGenerate(),
			cmdGet(),
			cmdInstall(),
			cmdList(),
			cmdMod(),
			cmdWork(),
			cmdRun(),
			cmdTelemetry(),
			cmdTest(),
			cmdTool(),
			cmdVersion(),
			cmdVet(),
		},
	}
}

func main() {
	ctx := context.Background()
	cmd := cmdGo()
	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
