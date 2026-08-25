package chvorinov

import (
	"math"
	"testing"
)

func TestSphereSlowestPlateFastest(t *testing.T) {
	table := CompareShapes(1000.0, 3.0, 2.0, 0)
	if !table.IsSphereSlowest() {
		t.Errorf("slowest = %s, want sphere", table.Slowest)
	}
	if !table.IsPlateFastest() {
		t.Errorf("fastest = %s, want plate", table.Fastest)
	}
	if len(table.Rows) != 4 {
		t.Errorf("len(rows) = %d, want 4", len(table.Rows))
	}
}

func TestSphereModulusLargestAmongShapes(t *testing.T) {
	table := CompareShapes(1000.0, 3.0, 2.0, 0)
	sphereM := 0.0
	for _, r := range table.Rows {
		if r.Label == "sphere" {
			sphereM = r.Modulus
		}
	}
	for _, r := range table.Rows {
		if r.Label == "sphere" {
			continue
		}
		if r.Modulus >= sphereM {
			t.Errorf("shape %s modulus %v should be below sphere modulus %v", r.Label, r.Modulus, sphereM)
		}
	}
}

func TestEqualModulusSameTime(t *testing.T) {
	m := 1.667
	t1 := EqualModulusTime(m, 3.0, 2.0)
	t2 := EqualModulusTime(m, 3.0, 2.0)
	if t1 != t2 {
		t.Errorf("times differ for equal moduli: %v vs %v", t1, t2)
	}
}

func TestFlatterPlateShorterTime(t *testing.T) {
	thick := 5.0
	flat := 2.5
	thickM := IdealPlateModulus(thick)
	flatM := IdealPlateModulus(flat)
	thickTf := FreezeTime(3.0, thickM, 2.0)
	flatTf := FreezeTime(3.0, flatM, 2.0)
	if flatTf >= thickTf {
		t.Errorf("flat tf %v should be below thick tf %v", flatTf, thickTf)
	}
}

func TestCubeModulusForVolumePinsReference(t *testing.T) {
	v := 1000.0
	got := CubeModulusForVolume(v)
	want := 10.0 / 6.0
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("CubeModulusForVolume(1000) = %v, want %v", got, want)
	}
}
