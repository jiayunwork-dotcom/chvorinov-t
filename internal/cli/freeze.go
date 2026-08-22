// This file implements the freeze subcommand, the main entry point of the
// tool. Given a casting JSON file it computes the modulus M = V/A and the
// freezing time tf = C * M^n, prints the physical feasibility check, the
// shape comparison table at the same volume, the optional riser check, and
// the similarity-scaling cross rule. A failed computation never prints a
// partial result: errors are returned to the router which reports them on
// stderr with a non-zero exit code.
package cli

import (
	"io"

	"chvorinov-t/internal/chvorinov"
)

// runFreeze executes the freeze subcommand for one JSON input file.
func runFreeze(stdout, stderr io.Writer, path string) int {
	in, err := LoadInput(path)
	if err != nil {
		return commandError(stderr, err)
	}

	shape, err := in.geometryShape()
	if err != nil {
		return commandError(stderr, err)
	}
	n := in.exponent()

	res, err := chvorinov.Compute(shape.V, shape.A, in.C, n, in.SuperheatK)
	if err != nil {
		return commandError(stderr, err)
	}

	writeResult(stdout, res, in.Label)

	table := chvorinov.CompareShapes(shape.V, in.C, n, 0)
	writeShapeTable(stdout, table)

	if in.Riser != nil {
		rm, ok, err := in.riserModulus()
		if err != nil {
			return commandError(stderr, err)
		}
		if ok {
			check := chvorinov.CheckRiser(res.Modulus, rm)
			writeRiser(stdout, check)
			if check.HasWarning() {
				// An undersized riser is a warning, not a failure: the
				// casting computation itself is still valid. The warning is
				// repeated on stderr so scripts that grep for it see it.
				writeWarning(stderr, check.Warning)
			}
		}
	}

	rep := chvorinov.CrossScale(res.Modulus, res.FreezeTime, 2.0, n)
	writeScale(stdout, rep)

	return ExitOK
}
