package chvorinov

import (
	"math"
	"testing"
)

func TestValidateRejectsZeroVolume(t *testing.T) {
	err := Validate(0.0, 600.0, 3.0, 2.0)
	if err == nil {
		t.Fatal("Validate(0, A, C, n) = nil, want error")
	}
}

func TestValidateRejectsNegativeArea(t *testing.T) {
	err := Validate(1000.0, -5.0, 3.0, 2.0)
	if err == nil {
		t.Fatal("Validate(V, -5, C, n) = nil, want error")
	}
}

func TestValidateRejectsZeroMoldConst(t *testing.T) {
	err := Validate(1000.0, 600.0, 0.0, 2.0)
	if err == nil {
		t.Fatal("Validate(V, A, 0, n) = nil, want error")
	}
}

func TestValidateRejectsNonPositiveExponent(t *testing.T) {
	for _, n := range []float64{0.0, -1.0, -0.5} {
		err := Validate(1000.0, 600.0, 3.0, n)
		if err == nil {
			t.Errorf("Validate(V, A, C, %v) = nil, want error", n)
		}
	}
}

func TestValidateRejectsNaN(t *testing.T) {
	if err := Validate(math.NaN(), 600.0, 3.0, 2.0); err == nil {
		t.Error("Validate(NaN, ...) = nil, want error")
	}
	if err := Validate(1000.0, 600.0, math.Inf(1), 2.0); err == nil {
		t.Error("Validate(V, A, Inf, n) = nil, want error")
	}
}

func TestValidateAcceptsPositiveInputs(t *testing.T) {
	if err := Validate(1000.0, 600.0, 3.0, 2.0); err != nil {
		t.Errorf("Validate(valid) = %v, want nil", err)
	}
}

func TestValidateInputRejectsNegativeSuperheat(t *testing.T) {
	in := Input{Volume: 1000, Area: 600, MoldConst: 3, Exponent: 2, SuperheatK: -10}
	if err := ValidateInput(in); err == nil {
		t.Error("ValidateInput(superheat=-10) = nil, want error")
	}
}

func TestValidationErrorNamesField(t *testing.T) {
	err := Validate(0.0, 600.0, 3.0, 2.0)
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T, want *ValidationError", err)
	}
	if ve.Field != "volume" {
		t.Errorf("field = %q, want \"volume\"", ve.Field)
	}
}
