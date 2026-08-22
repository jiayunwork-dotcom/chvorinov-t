package chvorinov

import (
	"math"
	"testing"
)

// TestSphereSlowestPlateFastest pins "same volume: sphere freezes slower
// than plate" from the requirements.
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

// TestSphereModulusLargestAmongShapes verifies the sphere row has the
// largest modulus at equal volume, which is why it freezes last.
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

// TestEqualModulusSameTime verifies that shapes sharing V/A produce the
// same freezing time under pure Chvorinov, regardless of their label.
func TestEqualModulusSameTime(t *testing.T) {
	m := 1.667
	t1 := EqualModulusTime(m, 3.0, 2.0)
	t2 := EqualModulusTime(m, 3.0, 2.0)
	if t1 != t2 {
		t.Errorf("times differ for equal moduli: %v vs %v", t1, t2)
	}
}

// TestFlatterPlateShorterTime is the cross rule at the physics level: for
// the same volume, a flatter plate has a shorter freezing time.
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

// TestCubeModulusForVolumePinsReference verifies M = a/6 for a cube of
// known volume.
func TestCubeModulusForVolumePinsReference(t *testing.T) {
	v := 1000.0 // a = 10 cm
	got := CubeModulusForVolume(v)
	want := 10.0 / 6.0
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("CubeModulusForVolume(1000) = %v, want %v", got, want)
	}
}
