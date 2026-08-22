package chvorinov

import "testing"

// TestRiserOk verifies that a riser whose modulus exceeds the casting
// modulus solidifies last and passes without a warning.
func TestRiserOk(t *testing.T) {
	check := CheckRiser(1.667, 2.0)
	if check.Status != RiserOk {
		t.Errorf("status = %v, want ok", check.Status)
	}
	if check.HasWarning() {
		t.Errorf("warning = %q, want none", check.Warning)
	}
}

// TestRiserTooSmall verifies that a riser with a smaller modulus than the
// casting produces the "must exceed" warning.
func TestRiserTooSmall(t *testing.T) {
	check := CheckRiser(1.667, 1.2)
	if check.Status != RiserTooSmall {
		t.Errorf("status = %v, want too-small", check.Status)
	}
	if !check.HasWarning() {
		t.Error("no warning for an undersized riser")
	}
}

// TestRiserEqualFlagsNoMargin verifies the equal-modulus case is called
// out as having no feeding margin.
func TestRiserEqualFlagsNoMargin(t *testing.T) {
	check := CheckRiser(1.667, 1.667)
	if check.Status != RiserEqual {
		t.Errorf("status = %v, want equal", check.Status)
	}
	if !check.HasWarning() {
		t.Error("no warning for an equal-modulus riser")
	}
}

// TestMinRiserModulusReturnsRequiredValue pins the required modulus.
func TestMinRiserModulusReturnsRequiredValue(t *testing.T) {
	got := MinRiserModulus(1.667)
	if got != 1.667 {
		t.Errorf("MinRiserModulus(1.667) = %v, want 1.667", got)
	}
}

// TestRiserDescribeRendersLine verifies the one-line summary has content.
func TestRiserDescribeRendersLine(t *testing.T) {
	check := CheckRiser(1.667, 2.0)
	if check.Describe() == "" {
		t.Fatal("Describe() = empty")
	}
}

// TestValidateRiserInputsRejectsNonPositive verifies the riser-side guard.
func TestValidateRiserInputsRejectsNonPositive(t *testing.T) {
	if err := ValidateRiserInputs(0.0); err == nil {
		t.Error("ValidateRiserInputs(0) = nil, want error")
	}
	if err := ValidateRiserInputs(2.0); err != nil {
		t.Errorf("ValidateRiserInputs(2) = %v, want nil", err)
	}
}
