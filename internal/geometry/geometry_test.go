package geometry

import (
	"math"
	"testing"
)

// TestCubeModulusAOver6 pins the fixed reference point of the whole tool:
// for a cube the modulus is exactly M = a/6.
func TestCubeModulusAOver6(t *testing.T) {
	cases := []struct {
		edge float64
		want float64
	}{
		{6.0, 1.0},
		{10.0, 10.0 / 6.0},
		{2.5, 2.5 / 6.0},
	}
	for _, c := range cases {
		got := CubeModulus(c.edge)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("CubeModulus(%v) = %v, want %v", c.edge, got, c.want)
		}
	}
}

// TestCubeVolumeArea verifies V = a^3 and A = 6a^2.
func TestCubeVolumeArea(t *testing.T) {
	a := 4.0
	wantV := 64.0
	wantA := 96.0
	if got := CubeVolume(a); got != wantV {
		t.Errorf("CubeVolume(4) = %v, want %v", got, wantV)
	}
	if got := CubeArea(a); got != wantA {
		t.Errorf("CubeArea(4) = %v, want %v", got, wantA)
	}
}

// TestSphereModulusROver3 pins the textbook sphere relation M = r/3.
func TestSphereModulusROver3(t *testing.T) {
	for _, r := range []float64{1.0, 3.0, 6.5} {
		got := SphereModulus(r)
		want := r / 3.0
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("SphereModulus(%v) = %v, want %v", r, got, want)
		}
	}
}

// TestPlateModulusApproachesHalfThickness checks that the finite plate
// modulus converges to t/2 as the plate becomes thinner relative to its
// face dimensions.
func TestPlateModulusApproachesHalfThickness(t *testing.T) {
	tol := 0.05 // 5% relative tolerance
	cases := []struct {
		t, w, l float64
	}{
		{1.0, 100.0, 100.0},
		{0.5, 200.0, 200.0},
		{2.0, 500.0, 100.0},
	}
	for _, c := range cases {
		got := PlateModulus(c.t, c.w, c.l)
		ideal := PlateModulusInfinite(c.t)
		if math.Abs(got-ideal)/ideal > tol {
			t.Errorf("PlateModulus(%v,%v,%v) = %v, want within 5%% of t/2 = %v", c.t, c.w, c.l, got, ideal)
		}
	}
}

// TestCylinderModulusApproachesHalfRadius checks convergence of the finite
// cylinder modulus toward r/2 for tall cylinders.
func TestCylinderModulusApproachesHalfRadius(t *testing.T) {
	tol := 0.05
	for _, r := range []float64{1.0, 2.0, 5.0} {
		h := 50.0 * r
		got := CylinderModulus(r, h)
		ideal := CylinderModulusInfinite(r)
		if math.Abs(got-ideal)/ideal > tol {
			t.Errorf("CylinderModulus(%v,%v) = %v, want within 5%% of r/2 = %v", r, h, got, ideal)
		}
	}
}

// TestShapeNewComputesGeometry verifies that New fills V and A from the
// shape kind and its dimensions.
func TestShapeNewComputesGeometry(t *testing.T) {
	s := New(ShapeCube, Dims{Edge: 2.0})
	if got, want := s.V, 8.0; got != want {
		t.Errorf("cube volume = %v, want %v", got, want)
	}
	if got, want := s.A, 24.0; got != want {
		t.Errorf("cube area = %v, want %v", got, want)
	}
	if got, want := s.Modulus(), 8.0/24.0; got != want {
		t.Errorf("cube modulus = %v, want %v", got, want)
	}
}

// TestFromVolumeAreaKeepsKindCustom verifies that custom shapes keep their
// label and use the raw volume and area.
func TestFromVolumeAreaKeepsKindCustom(t *testing.T) {
	s := FromVolumeArea(10.0, 20.0)
	if s.Kind != ShapeCustom {
		t.Errorf("kind = %v, want custom", s.Kind)
	}
	if got, want := s.Modulus(), 0.5; got != want {
		t.Errorf("modulus = %v, want %v", got, want)
	}
}

// TestMinSurfaceAreaIsEnclosingSphere checks the isoperimetric bound.
func TestMinSurfaceAreaIsEnclosingSphere(t *testing.T) {
	v := 1000.0
	r := RadiusFromVolume(v)
	want := 4 * math.Pi * r * r
	got := MinSurfaceArea(v)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("MinSurfaceArea(%v) = %v, want %v", v, got, want)
	}
}

// TestNearlyEqualReportsTolerance verifies the shared tolerance helper.
func TestNearlyEqualReportsTolerance(t *testing.T) {
	if !NearlyEqual(1.0, 1.0+1e-12) {
		t.Error("NearlyEqual(1, 1+1e-12) = false, want true")
	}
	if NearlyEqual(1.0, 2.0) {
		t.Error("NearlyEqual(1, 2) = true, want false")
	}
}
