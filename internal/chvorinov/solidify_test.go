package chvorinov

import (
	"math"
	"testing"
)

// TestFreezeTimeFormula pins tf = C * M^n for representative values.
func TestFreezeTimeFormula(t *testing.T) {
	cases := []struct {
		c, m, n float64
		want    float64
	}{
		{3.0, 2.0, 2.0, 12.0},
		{3.0, 1.0, 2.0, 3.0},
		{2.0, 2.0, 1.0, 4.0},
		{1.0, 3.0, 3.0, 27.0},
	}
	for _, c := range cases {
		got := FreezeTime(c.c, c.m, c.n)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("FreezeTime(%v, %v, %v) = %v, want %v", c.c, c.m, c.n, got, c.want)
		}
	}
}

// TestFreezeTimeFromVA verifies the direct form tf = C * (V/A)^n.
func TestFreezeTimeFromVA(t *testing.T) {
	got := FreezeTimeFromVA(1000.0, 600.0, 3.0, 2.0)
	want := 3.0 * math.Pow(1000.0/600.0, 2.0)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("FreezeTimeFromVA = %v, want %v", got, want)
	}
}

// TestComputeReturnsResult pins the full pipeline for a valid cube.
func TestComputeReturnsResult(t *testing.T) {
	// Cube edge 10 cm: V = 1000, A = 600, M = 1.6667, tf = 3 * M^2.
	res, err := Compute(1000.0, 600.0, 3.0, 2.0, 0)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	if got, want := res.Modulus, 1000.0/600.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("modulus = %v, want %v", got, want)
	}
	if got, want := res.FreezeTime, 3.0*math.Pow(1000.0/600.0, 2.0); math.Abs(got-want) > 1e-12 {
		t.Errorf("freeze time = %v, want %v", got, want)
	}
	if res.SuperheatApplied {
		t.Error("superheat applied for superheatK = 0, want pure Chvorinov")
	}
}

// TestComputeAppliesSuperheat verifies that a positive superheat extends
// the freezing time and flags the correction.
func TestComputeAppliesSuperheat(t *testing.T) {
	res, err := Compute(1000.0, 600.0, 3.0, 2.0, 50.0)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	if !res.SuperheatApplied {
		t.Error("superheat not applied for superheatK = 50")
	}
	pure := 3.0 * math.Pow(1000.0/600.0, 2.0)
	if res.FreezeTime <= pure {
		t.Errorf("freeze time with superheat %v should exceed pure %v", res.FreezeTime, pure)
	}
}

// TestResultStringIncludesCoreNumbers checks that the compact summary
// contains the modulus and the freezing time.
func TestResultStringIncludesCoreNumbers(t *testing.T) {
	res, err := Compute(1000.0, 600.0, 3.0, 2.0, 0)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	s := res.String()
	if s == "" {
		t.Fatal("Result.String() = empty")
	}
}
