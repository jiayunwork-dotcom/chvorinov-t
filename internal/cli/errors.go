package cli

import (
	"fmt"
	"io"
)

const (
	ExitOK    = 0
	ExitError = 1
	ExitUsage = 2
)

var Version = "0.1.0"

func writeError(w io.Writer, err error) {
	fmt.Fprintf(w, "chvorinov-t: error: %v\n", err)
}

func usageError(w io.Writer, msg string) int {
	fmt.Fprintf(w, "chvorinov-t: %s\n", msg)
	fmt.Fprintf(w, "run \"chvorinov-t --help\" for usage\n")
	return ExitUsage
}

func commandError(w io.Writer, err error) int {
	writeError(w, err)
	return ExitError
}

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

func errNonPositive(field string, value float64) error {
	return fmt.Errorf("%s must be strictly positive (got %v)", field, value)
}

var errNoRiser = fmt.Errorf("input file has no \"riser\" section")

func fmtLine(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format+"\n", args...)
}
