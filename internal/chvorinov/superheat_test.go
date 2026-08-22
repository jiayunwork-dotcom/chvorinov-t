package chvorinov

import (
	"math"
	"testing"
)

// TestSuperheatFactorZeroIsPure verifies that zero superheat leaves the
// pure Chvorinov time untouched (factor 1).
func TestSuperheatFactorZeroIsPure(t *testing.T) {
	got := SuperheatFactor(0.0, SteelHeatCapacity, SteelLatentHeat)
	if got != 1.0 {
		t.Errorf("SuperheatFactor(0) = %v, want 1", got)
	}
}

// TestSuperheatFactorPositive verifies the 1 + cp*dT/Lf form with the
// package steel constants: 50 K yields roughly 1.148.
func TestSuperheatFactorPositive(t *testing.T) {
	got := SuperheatFactor(50.0, SteelHeatCapacity, SteelLatentHeat)
	want := 1.0 + 0.80*50.0/270.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("SuperheatFactor(50) = %v, want %v", got, want)
	}
}

// TestApplySuperheatMultipliesTime pins the extended freezing time.
func TestApplySuperheatMultipliesTime(t *testing.T) {
	got := ApplySuperheat(8.333, 50.0, SteelHeatCapacity, SteelLatentHeat)
	want := 8.333 * (1.0 + 0.80*50.0/270.0)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("ApplySuperheat = %v, want %v", got, want)
	}
}

// TestSuperheatExtensionIsNonNegative verifies the extension is zero for
// pure Chvorinov and positive for real superheat.
func TestSuperheatExtensionIsNonNegative(t *testing.T) {
	if got := SuperheatExtension(8.333, 0.0, SteelHeatCapacity, SteelLatentHeat); got != 0.0 {
		t.Errorf("extension at 0 superheat = %v, want 0", got)
	}
	if got := SuperheatExtension(8.333, 50.0, SteelHeatCapacity, SteelLatentHeat); got <= 0 {
		t.Errorf("extension at 50 K = %v, want > 0", got)
	}
}

// TestSteelSuperheatFactorMatchesGeneric verifies the convenience wrapper.
func TestSteelSuperheatFactorMatchesGeneric(t *testing.T) {
	generic := SuperheatFactor(100.0, SteelHeatCapacity, SteelLatentHeat)
	if got := SteelSuperheatFactor(100.0); got != generic {
		t.Errorf("SteelSuperheatFactor = %v, want %v", got, generic)
	}
}
