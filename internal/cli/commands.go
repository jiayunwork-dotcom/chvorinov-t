// This file implements the remaining subcommands. compare prints the
// per-shape modulus table at the volume of the input file, scale verifies
// the similarity-scaling cross rules, riser checks an attached riser in
// isolation, and validate runs the full input validation without printing
// any computed result. Each command is deliberately small because the
// heavy lifting lives in the physics packages; the bodies here only move
// data between the input document and the report writers.
package cli

import (
	"io"

	"chvorinov-t/internal/chvorinov"
)

// runCompare executes the compare subcommand: a shape comparison table at
// the volume described by the input file. The shape class of the file is
// irrelevant here; only its volume, mold constant, and exponent matter.
func runCompare(stdout, stderr io.Writer, path string) int {
	in, err := LoadInput(path)
	if err != nil {
		return commandError(stderr, err)
	}

	shape, err := in.geometryShape()
	if err != nil {
		return commandError(stderr, err)
	}
	if shape.V <= 0 {
		return commandError(stderr, errNonPositive("volume", shape.V))
	}
	n := in.exponent()

	table := chvorinov.CompareShapes(shape.V, in.C, n, 0)
	writeShapeTable(stdout, table)
	return ExitOK
}

// runScale executes the scale subcommand: it verifies that scaling every
// linear dimension by two doubles the modulus and, with n = 2, quadruples
// the freezing time, and that doubling the mold constant doubles the time.
func runScale(stdout, stderr io.Writer, path string) int {
	in, err := LoadInput(path)
	if err != nil {
		return commandError(stderr, err)
	}

	shape, err := in.geometryShape()
	if err != nil {
		return commandError(stderr, err)
	}
	n := in.exponent()

	res, err := chvorinov.Compute(shape.V, shape.A, in.C, n, 0)
	if err != nil {
		return commandError(stderr, err)
	}

	rep := chvorinov.CrossScale(res.Modulus, res.FreezeTime, 2.0, n)
	writeScale(stdout, rep)

	writeHeading(stdout, "mold constant doubling")
	td := chvorinov.DoubleMoldConstTime(res.FreezeTime)
	fmtLine(stdout, "tf(C)  : %s min", fmtNum(res.FreezeTime))
	fmtLine(stdout, "tf(2C) : %s min  (tf x2)", fmtNum(td))

	writeHeading(stdout, "equal-modulus cross rule")
	fmtLine(stdout, "any shape with M=%.4f cm has tf=%.3f min", res.Modulus, res.FreezeTime)
	return ExitOK
}

// runRiser executes the riser subcommand: it computes the casting modulus
// and the riser modulus from the input file and reports whether the riser
// can feed the casting (i.e. whether it solidifies last).
func runRiser(stdout, stderr io.Writer, path string) int {
	in, err := LoadInput(path)
	if err != nil {
		return commandError(stderr, err)
	}

	shape, err := in.geometryShape()
	if err != nil {
		return commandError(stderr, err)
	}
	n := in.exponent()

	res, err := chvorinov.Compute(shape.V, shape.A, in.C, n, 0)
	if err != nil {
		return commandError(stderr, err)
	}

	writeResult(stdout, res, in.Label)

	rm, ok, err := in.riserModulus()
	if err != nil {
		return commandError(stderr, err)
	}
	if !ok {
		return commandError(stderr, errNoRiser)
	}

	check := chvorinov.CheckRiser(res.Modulus, rm)
	writeRiser(stdout, check)
	if check.HasWarning() {
		writeWarning(stderr, check.Warning)
	}
	return ExitOK
}

// runValidate executes the validate subcommand: it loads the input file
// and runs the same validation pipeline as freeze, printing only a
// success line when everything passes.
func runValidate(stdout, stderr io.Writer, path string) int {
	in, err := LoadInput(path)
	if err != nil {
		return commandError(stderr, err)
	}

	shape, err := in.geometryShape()
	if err != nil {
		return commandError(stderr, err)
	}
	n := in.exponent()

	if _, err := chvorinov.Compute(shape.V, shape.A, in.C, n, in.SuperheatK); err != nil {
		return commandError(stderr, err)
	}
	fmtLine(stdout, "valid: %s", path)
	return ExitOK
}
