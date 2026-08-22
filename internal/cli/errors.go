// This file centralises the CLI error policy. Every failure path prints a
// single "error:" line to stderr and exits with a non-zero status; no
// error is ever swallowed or downgraded to a warning in the main path.
// The exit codes distinguish usage mistakes (2) from invalid inputs and
// runtime failures (1) so scripts can tell them apart.
package cli

import (
	"fmt"
	"io"
)

// Exit codes returned by Run.
const (
	// ExitOK means the command completed successfully.
	ExitOK = 0
	// ExitError means the command failed on invalid input or a runtime
	// problem such as a missing file.
	ExitError = 1
	// ExitUsage means the command line itself was malformed: unknown
	// subcommand or missing arguments.
	ExitUsage = 2
)

// Version is the release string printed by the version subcommand. It is a
// variable so tests can override it, and the build can inject a tag.
var Version = "0.1.0"

// writeError prints one error line to stderr in the conventional
// "chvorinov-t: message" form used by Unix command-line tools.
func writeError(w io.Writer, err error) {
	fmt.Fprintf(w, "chvorinov-t: error: %v\n", err)
}

// usageError reports a command-line usage problem together with a hint to
// run --help, and returns the usage exit code.
func usageError(w io.Writer, msg string) int {
	fmt.Fprintf(w, "chvorinov-t: %s\n", msg)
	fmt.Fprintf(w, "run \"chvorinov-t --help\" for usage\n")
	return ExitUsage
}

// commandError reports an invalid-input or runtime failure and returns the
// error exit code.
func commandError(w io.Writer, err error) int {
	writeError(w, err)
	return ExitError
}

// printUsage writes the full usage text. It is kept here so both --help
// and the unknown-subcommand path render the same text.
func printUsage(w io.Writer) {
	fmt.Fprint(w, `chvorinov-t: casting solidification time calculator (Chvorinov's rule)

Usage:
  chvorinov-t freeze <file.json>    compute M and tf; print shape comparison
  chvorinov-t compare <file.json>   shape modulus comparison at fixed volume
  chvorinov-t scale <file.json>     similarity-scaling cross rule (size x2)
  chvorinov-t riser <file.json>     riser modulus check for a casting
  chvorinov-t validate <file.json>  validate a casting file without computing
  chvorinov-t version               print the version
  chvorinov-t --help                show this help

The JSON file describes a casting with shape, dimensions (or volume/area),
the mold constant C, and an optional exponent n (default 2). Units are
centimetre-based: V in cm^3, A in cm^2, M in cm, C in min/cm^n, tf in min.

Example:
  chvorinov-t freeze example/steel-cube.json
`)
}

// errNonPositive builds a simple "field must be positive" error used by the
// subcommands when a derived quantity is invalid before the physics layer
// runs.
func errNonPositive(field string, value float64) error {
	return fmt.Errorf("%s must be strictly positive (got %v)", field, value)
}

// errNoRiser is returned when a command that requires a riser is given an
// input file without one.
var errNoRiser = fmt.Errorf("input file has no \"riser\" section")

// fmtLine writes a formatted line terminated by a newline to a writer. It
// is the report equivalent of fmt.Fprintf with a trailing "\n".
func fmtLine(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format+"\n", args...)
}
