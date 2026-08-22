package chvorinov

import (
	"math"
	"testing"
)

// TestScaleDoublingQuadruplesTime is the canonical cross rule from the
// requirements: doubling every linear dimension doubles M and, at n = 2,
// quadruples the freezing time.
func TestScaleDoublingQuadruplesTime(t *testing.T) {
	m := 2.0
	tf := 3.0
	rep := CrossScale(m, tf, 2.0, 2.0)
	if got, want := rep.ScaledModulus, 4.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("scaled modulus = %v, want %v", got, want)
	}
	if got, want := rep.ScaledTime, 12.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("scaled time = %v, want %v", got, want)
	}
	if !IsQuadrupleTime(rep.TimeRatio) {
		t.Errorf("time ratio = %v, want 4", rep.TimeRatio)
	}
}

// TestScaleModulusScalesLinearly verifies M' = 2M under similarity scaling.
func TestScaleModulusScalesLinearly(t *testing.T) {
	got := ScaleModulus(1.6667, 2.0)
	want := 3.3334
	if math.Abs(got-want) > 1e-3 {
		t.Errorf("ScaleModulus(1.6667, 2) = %v, want %v", got, want)
	}
}

// TestCDoublingDoublesTime pins the "C x2 => tf x2" cross rule.
func TestCDoublingDoublesTime(t *testing.T) {
	tf := 8.333
	got := DoubleMoldConstTime(tf)
	want := 16.666
	if math.Abs(got-want) > 1e-3 {
		t.Errorf("DoubleMoldConstTime(%v) = %v, want %v", tf, got, want)
	}
	if !IsDoubleTime(got / tf) {
		t.Errorf("time ratio = %v, want 2", got/tf)
	}
}

// TestScaleFactorFromTimesInvertsScaleTime verifies that recovering the
// scale factor from a time ratio works both ways.
func TestScaleFactorFromTimesInvertsScaleTime(t *testing.T) {
	factor := ScaleFactorFromTimes(8.333, 33.332, 2.0)
	if math.Abs(factor-2.0) > 1e-3 {
		t.Errorf("recovered factor = %v, want 2.0", factor)
	}
}

// TestVerifyCrossRulesAllPass checks every cross rule in one call.
func TestVerifyCrossRulesAllPass(t *testing.T) {
	cr := VerifyCrossRules(1000.0, 3.0, 2.0)
	if !cr.AllPass {
		for _, r := range cr.Rules {
			t.Errorf("rule %q: pass=%v detail=%s", r.Name, r.Pass, r.Detail)
		}
		t.Fatal("VerifyCrossRules reported a failing rule")
	}
	if len(cr.Rules) < 5 {
		t.Errorf("len(rules) = %d, want at least 5", len(cr.Rules))
	}
}
