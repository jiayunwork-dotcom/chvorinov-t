package chvorinov

import "testing"

func TestRiserOk(t *testing.T) {
	check := CheckRiser(1.667, 2.0)
	if check.Status != RiserOk {
		t.Errorf("status = %v, want ok", check.Status)
	}
	if check.HasWarning() {
		t.Errorf("warning = %q, want none", check.Warning)
	}
}

func TestRiserTooSmall(t *testing.T) {
	check := CheckRiser(1.667, 1.2)
	if check.Status != RiserTooSmall {
		t.Errorf("status = %v, want too-small", check.Status)
	}
	if !check.HasWarning() {
		t.Error("no warning for an undersized riser")
	}
}

func TestRiserEqualFlagsNoMargin(t *testing.T) {
	check := CheckRiser(1.667, 1.667)
	if check.Status != RiserEqual {
		t.Errorf("status = %v, want equal", check.Status)
	}
	if !check.HasWarning() {
		t.Error("no warning for an equal-modulus riser")
	}
}

func TestMinRiserModulusReturnsRequiredValue(t *testing.T) {
	got := MinRiserModulus(1.667)
	if got != 1.667 {
		t.Errorf("MinRiserModulus(1.667) = %v, want 1.667", got)
	}
}

func TestRiserDescribeRendersLine(t *testing.T) {
	check := CheckRiser(1.667, 2.0)
	if check.Describe() == "" {
		t.Fatal("Describe() = empty")
	}
}

func TestValidateRiserInputsRejectsNonPositive(t *testing.T) {
	if err := ValidateRiserInputs(0.0); err == nil {
		t.Error("ValidateRiserInputs(0) = nil, want error")
	}
	if err := ValidateRiserInputs(2.0); err != nil {
		t.Errorf("ValidateRiserInputs(2) = %v, want nil", err)
	}
}
