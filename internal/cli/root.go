package cli

import (
	"io"
)

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

type commandFunc func(stdout, stderr io.Writer, path string) int

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
