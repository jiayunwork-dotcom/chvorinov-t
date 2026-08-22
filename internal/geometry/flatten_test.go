package geometry

import "testing"

// TestFlatterPlateLargerAreaSmallerModulus pins the cross rule "same V,
// flatter plate has larger A and smaller M" at the geometry level.
func TestFlatterPlateLargerAreaSmallerModulus(t *testing.T) {
	v := 1000.0
	eff := Flatten(v, 10.0, 2.5)
	if !eff.IsFlatterLargerArea() {
		t.Errorf("flatter plate area = %v, want > ref area %v", eff.FlatArea, eff.RefArea)
	}
	if !eff.IsFlatterSmallerModulus() {
		t.Errorf("flatter plate modulus = %v, want < ref modulus %v", eff.FlatModulus, eff.RefModulus)
	}
	if !eff.IsFlatterShorterTime() {
		t.Errorf("flatter plate should have a shorter freezing time (monotone in M)")
	}
}

// TestPlateSeriesForVolumeBuildsOrderedCases verifies that thinner plates
// in a series have larger areas and smaller moduli, in order.
func TestPlateSeriesForVolumeBuildsOrderedCases(t *testing.T) {
	v := 1000.0
	series := PlateSeriesForVolume(v, []float64{10.0, 5.0, 2.5})
	if len(series) != 3 {
		t.Fatalf("len(series) = %d, want 3", len(series))
	}
	for i := 1; i < len(series); i++ {
		if series[i].Area <= series[i-1].Area {
			t.Errorf("case %d area %v should exceed case %d area %v", i, series[i].Area, i-1, series[i-1].Area)
		}
		if series[i].Modulus >= series[i-1].Modulus {
			t.Errorf("case %d modulus %v should be below case %d modulus %v", i, series[i].Modulus, i-1, series[i-1].Modulus)
		}
	}
}

// TestFlattenThinnerThicknessYieldsShorterTimeRatio verifies that the time
// ratio implied by the modulus ratio is below one, using a quadratic
// exponent.
func TestFlattenThinnerThicknessYieldsShorterTimeRatio(t *testing.T) {
	eff := Flatten(1000.0, 10.0, 5.0)
	ratio := eff.FlatModulus / eff.RefModulus // = (tf2/tf1)^(1/2) for n=2
	if ratio >= 1.0 {
		t.Errorf("modulus ratio = %v, want < 1 for a thinner plate", ratio)
	}
}
