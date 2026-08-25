package cli

import (
	"io"

	"chvorinov-t/internal/chvorinov"
)

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
