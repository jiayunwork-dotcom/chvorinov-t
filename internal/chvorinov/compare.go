package chvorinov

import (
	"fmt"
	"math"

	"chvorinov-t/internal/geometry"
)

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

type ShapeTable struct {
	Volume    float64
	MoldConst float64
	Exponent  float64
	Rows      []ShapeRow
	Slowest   string
	Fastest   string
}

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

func newRow(label, sizeDesc string, s geometry.Shape, c, n float64, formula string) ShapeRow {
	m := s.Modulus()
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

func timeOf(rows []ShapeRow, label string) float64 {
	for _, r := range rows {
		if r.Label == label {
			return r.Time
		}
	}
	return 0
}

func EqualModulusTime(m, moldConst, exponent float64) float64 {
	return FreezeTime(moldConst, m, exponent)
}

func (t ShapeTable) IsSphereSlowest() bool {
	return t.Slowest == "sphere"
}

func (t ShapeTable) IsPlateFastest() bool {
	return t.Fastest == "plate"
}
