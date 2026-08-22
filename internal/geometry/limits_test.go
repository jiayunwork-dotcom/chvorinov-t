package geometry

import (
	"strings"
	"testing"
)

// TestAreaBelowMinimumRejected verifies that a surface area smaller than
// the enclosing sphere's area for the same volume is rejected.
func TestAreaBelowMinimumRejected(t *testing.T) {
	v := 1000.0
	a := MinSurfaceArea(v) * 0.5
	err := CheckFeasibility(v, a)
	if err == nil {
		t.Fatalf("CheckFeasibility(%v, %v) = nil, want error", v, a)
	}
	if !strings.Contains(err.Error(), "enclosing sphere") {
		t.Errorf("error %q does not mention the enclosing-sphere bound", err)
	}
}

// TestAreaAtMinimumAccepted verifies that exactly the sphere's area passes
// the feasibility check (floating-point noise tolerated).
func TestAreaAtMinimumAccepted(t *testing.T) {
	v := 1000.0
	a := MinSurfaceArea(v)
	if err := CheckFeasibility(v, a); err != nil {
		t.Errorf("CheckFeasibility(sphere area) = %v, want nil", err)
	}
}

// TestIsAreaBelowMinimum reports the boolean form of the bound.
func TestIsAreaBelowMinimum(t *testing.T) {
	v := 64.0
	minA := MinSurfaceArea(v)
	if !IsAreaBelowMinimum(v, minA*0.99) {
		t.Error("IsAreaBelowMinimum(0.99*min) = false, want true")
	}
	if IsAreaBelowMinimum(v, minA*1.01) {
		t.Error("IsAreaBelowMinimum(1.01*min) = true, want false")
	}
}

// TestAreaMarginFraction is the signed margin above the physical minimum.
func TestAreaMarginFraction(t *testing.T) {
	v := 64.0
	a := 2.0 * MinSurfaceArea(v)
	got := AreaMarginFraction(v, a)
	if got != 1.0 {
		t.Errorf("AreaMarginFraction(2*min) = %v, want 1.0", got)
	}
}
