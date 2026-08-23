// This file implements the shape comparison table. For a fixed volume the
// sphere exposes the smallest surface area, hence the largest modulus and
// the longest freezing time; a thin plate exposes much more surface and
// freezes faster. The table lists cube, sphere, thin plate, and a tall
// cylinder with their representative dimensions, exact volume-to-area
// ratios, and freezing times, so the "sphere slower than plate" and "a
// flatter plate at equal volume has a shorter tf" cross rules can be read
// off directly. The ideal-modulus formulas (plate t/2, cylinder r/2,
// sphere r/3, cube a/6) are reported per row as the textbook reference.
package chvorinov

import (
	"fmt"
	"math"

	"chvorinov-t/internal/geometry"
)

// ShapeRow is one line of the comparison table.
type ShapeRow struct {
	Label        string
	SizeDesc     string
	Volume       float64
	Area         float64
	Modulus      float64
	Time         float64
	IdealModulus float64
	IdealFormula string
}

// ShapeTable is the full comparison for one volume.
type ShapeTable struct {
	Volume    float64
	MoldConst float64
	Exponent  float64
	Rows      []ShapeRow
	Slowest   string
	Fastest   string
}

// representativeSizes derives a representative dimension set for each
// shape so that every row encloses exactly the same volume.
func representativeSizes(volume, plateThickness float64) (cubeA, sphereR, plateT, plateW, plateL, cylR, cylH float64) {
	cubeA = math.Cbrt(volume)
	sphereR = geometry.RadiusFromVolume(volume)

	plateT = cubeA / 4.0
	if plateThickness > 0 {
		plateT = plateThickness
	}
	plateW = math.Sqrt(volume / plateT)
	plateL = plateW

	cylR = math.Cbrt(volume / (5.0 * math.Pi))
	cylH = 5.0 * cylR
	return cubeA, sphereR, plateT, plateW, plateL, cylR, cylH
}

// CompareShapes builds the comparison table for a volume, mold constant,
// and exponent. An optional plateThickness overrides the default thin
// plate (cube edge / 4); pass 0 to keep the default. The volume must be
// positive; validation is the caller's responsibility.
func CompareShapes(volume, moldConst, exponent, plateThickness float64) ShapeTable {
	cubeA, sphereR, plateT, plateW, plateL, cylR, cylH := representativeSizes(volume, plateThickness)

	cube := geometry.Shape{
		Kind: geometry.ShapeCube,
		Dims: geometry.Dims{Edge: cubeA},
		V:    geometry.CubeVolume(cubeA),
		A:    geometry.CubeArea(cubeA),
	}
	sphere := geometry.Shape{
		Kind: geometry.ShapeSphere,
		Dims: geometry.Dims{Radius: sphereR},
		V:    geometry.SphereVolume(sphereR),
		A:    geometry.SphereArea(sphereR),
	}
	plate := geometry.Shape{
		Kind: geometry.ShapePlate,
		Dims: geometry.Dims{Thickness: plateT, Width: plateW, Length: plateL},
		V:    geometry.PlateVolume(plateT, plateW, plateL),
		A:    geometry.PlateArea(plateT, plateW, plateL),
	}
	cylinder := geometry.Shape{
		Kind: geometry.ShapeCylinder,
		Dims: geometry.Dims{Radius: cylR, Height: cylH},
		V:    geometry.CylinderVolume(cylR, cylH),
		A:    geometry.CylinderArea(cylR, cylH),
	}

	rows := []ShapeRow{
		newRow("cube", fmt.Sprintf("a=%.3f", cubeA), cube, moldConst, exponent, "a/6"),
		newRow("sphere", fmt.Sprintf("r=%.3f", sphereR), sphere, moldConst, exponent, "r/3"),
		newRow("plate", fmt.Sprintf("t=%.3f", plateT), plate, moldConst, exponent, "t/2"),
		newRow("cylinder", fmt.Sprintf("r=%.3f h=%.3f", cylR, cylH), cylinder, moldConst, exponent, "r/2"),
	}

	slowest, fastest := rows[0].Label, rows[0].Label
	for _, r := range rows {
		if r.Time > timeOf(rows, slowest) {
			slowest = r.Label
		}
		if r.Time < timeOf(rows, fastest) {
			fastest = r.Label
		}
	}

	return ShapeTable{
		Volume:    volume,
		MoldConst: moldConst,
		Exponent:  exponent,
		Rows:      rows,
		Slowest:   slowest,
		Fastest:   fastest,
	}
}

// newRow fills one comparison row from a geometry.Shape.
func newRow(label, sizeDesc string, s geometry.Shape, c, n float64, formula string) ShapeRow {
	m := shapeModCache.modulusFor(s.V, s.Modulus())
	return ShapeRow{
		Label:        label,
		SizeDesc:     sizeDesc,
		Volume:       s.V,
		Area:         s.A,
		Modulus:      m,
		Time:         FreezeTime(c, m, n),
		IdealModulus: geometry.IdealModulus(s.Kind, s.Dims),
		IdealFormula: formula,
	}
}

// timeOf looks up the freezing time of a row by label. It is a small
// helper used while computing the slowest and fastest rows.
func timeOf(rows []ShapeRow, label string) float64 {
	for _, r := range rows {
		if r.Label == label {
			return r.Time
		}
	}
	return 0
}

// EqualModulusTime returns the freezing time of a body with a given
// modulus. Because tf depends only on M under pure Chvorinov, any two
// shapes that share V/A also share tf regardless of their shape labels.
func EqualModulusTime(m, moldConst, exponent float64) float64 {
	return FreezeTime(moldConst, m, exponent)
}

// IsSphereSlowest reports whether the sphere row has the longest freezing
// time in the table, pinning the "same volume sphere freezes slower than
// plate" cross rule.
func (t ShapeTable) IsSphereSlowest() bool {
	return t.Slowest == "sphere"
}

// IsPlateFastest reports whether the plate row has the shortest time.
func (t ShapeTable) IsPlateFastest() bool {
	return t.Fastest == "plate"
}
