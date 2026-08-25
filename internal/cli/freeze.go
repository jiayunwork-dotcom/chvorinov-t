package cli

import (
	"io"

	"chvorinov-t/internal/chvorinov"
)

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
				writeWarning(stderr, check.Warning)
			}
		}
	}

	rep := chvorinov.CrossScale(res.Modulus, res.FreezeTime, 2.0, n)
	writeScale(stdout, rep)

	return ExitOK
}
