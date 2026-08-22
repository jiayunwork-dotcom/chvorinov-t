// This file is the subcommand router. Run is the only entry point used by
// main; it dispatches on the first argument, binds the standard output and
// error writers, and returns the process exit code. Keeping the router
// separate from the command bodies makes the CLI testable without
// subprocesses.
package cli

import (
	"io"
)

// Run executes the chvorinov-t command line and returns the process exit
// code. args excludes the program name. stdout receives normal output,
// stderr receives errors and usage text.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return ExitUsage
	}

	switch args[0] {
	case "--help", "-h", "help":
		printUsage(stdout)
		return ExitOK
	case "version":
		writeVersion(stdout, Version)
		return ExitOK
	case "freeze":
		return requireFile(stdout, stderr, args[1:], runFreeze)
	case "compare":
		return requireFile(stdout, stderr, args[1:], runCompare)
	case "scale":
		return requireFile(stdout, stderr, args[1:], runScale)
	case "riser":
		return requireFile(stdout, stderr, args[1:], runRiser)
	case "validate":
		return requireFile(stdout, stderr, args[1:], runValidate)
	default:
		return usageError(stderr, "unknown command "+args[0])
	}
}

// commandFunc is the signature of a subcommand body: it receives the
// writers and the input file path.
type commandFunc func(stdout, stderr io.Writer, path string) int

// requireFile checks that exactly one file argument is present and then
// delegates to the command body.
func requireFile(stdout, stderr io.Writer, rest []string, body commandFunc) int {
	if len(rest) != 1 {
		return usageError(stderr, "expected exactly one JSON input file")
	}
	if rest[0] == "-h" || rest[0] == "--help" {
		printUsage(stdout)
		return ExitOK
	}
	return body(stdout, stderr, rest[0])
}
