// Command chvorinov-t is a command-line tool that computes casting
// solidification times with Chvorinov's rule. Given a casting volume V,
// surface area A, mold constant C, and exponent n, it computes the
// solidification modulus M = V/A and the freezing time tf = C * M^n, and
// compares the moduli of different shape classes (cube, plate, cylinder,
// sphere) at a fixed volume. It also checks attached risers and verifies
// the similarity-scaling cross rules.
//
// Examples:
//
//	go run . freeze example/steel-cube.json
//	go run . compare example/steel-cube.json
//	go run . version
//
// All physics lives in the internal/geometry and internal/chvorinov
// packages; this file only wires the command line to the subcommand
// router and converts the exit code into a process status.
package main

import (
	"os"

	"chvorinov-t/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
