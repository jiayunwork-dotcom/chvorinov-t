// This file holds the input validation for the Chvorinov kernel. The
// mathematical model requires V > 0, A > 0, C > 0, and n > 0; every other
// combination is rejected with a specific error so that a user who feeds
// an invalid JSON file sees exactly which quantity is wrong, instead of a
// division-by-zero panic or a silently meaningless result.
package chvorinov

import (
	"errors"
	"fmt"
	"math"
)

// Input is the validated physical description of one casting.
type Input struct {
	Volume     float64 // cm^3
	Area       float64 // cm^2
	MoldConst  float64 // min/cm^n
	Exponent   float64
	SuperheatK float64 // optional, zero means pure Chvorinov
}

// ValidationError describes a single invalid input field.
type ValidationError struct {
	Field string
	Value float64
	Rule  string
}

// Error implements the error interface with a message shaped for the CLI:
//
//	invalid volume: got 0, want > 0
func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: got %v, want %s", e.Field, e.Value, e.Rule)
}

// Sentinel errors so callers can test without string matching.
var (
	ErrVolumeNotPositive    = errors.New("volume must be strictly positive")
	ErrAreaNotPositive      = errors.New("area must be strictly positive")
	ErrMoldConstNotPositive = errors.New("mold constant must be strictly positive")
	ErrExponentNotPositive  = errors.New("exponent must be strictly positive")
	ErrNotFinite            = errors.New("value must be a finite number")
)

// IsNotFinite reports whether a value is NaN or infinite. JSON does not
// carry NaN or Inf directly, but a user can supply them through arithmetic
// in other tooling, so the validator defends against them explicitly.
func IsNotFinite(v float64) bool {
	return math.IsNaN(v) || math.IsInf(v, 0)
}

// checkFinite rejects NaN and infinite values for a named field.
func checkFinite(field string, v float64) error {
	if IsNotFinite(v) {
		return &ValidationError{Field: field, Value: v, Rule: "a finite number"}
	}
	return nil
}

// Validate checks the four physical inputs of the Chvorinov model and
// returns an error naming the first offending field. The checks are
// ordered so that the most fundamental quantity (volume) is reported
// first when several inputs are wrong at once.
func Validate(volume, area, moldConst, exponent float64) error {
	if err := checkFinite("volume", volume); err != nil {
		return err
	}
	if err := checkFinite("area", area); err != nil {
		return err
	}
	if err := checkFinite("mold constant", moldConst); err != nil {
		return err
	}
	if err := checkFinite("exponent", exponent); err != nil {
		return err
	}
	if volume <= 0 {
		return &ValidationError{Field: "volume", Value: volume, Rule: "> 0"}
	}
	if area <= 0 {
		return &ValidationError{Field: "area", Value: area, Rule: "> 0"}
	}
	if moldConst <= 0 {
		return &ValidationError{Field: "mold constant", Value: moldConst, Rule: "> 0"}
	}
	if exponent <= 0 {
		return &ValidationError{Field: "exponent", Value: exponent, Rule: "> 0"}
	}
	return nil
}

// ValidateInput runs Validate against a full Input and additionally
// verifies that the superheat value is non-negative. A negative superheat
// would imply the metal enters the mold colder than the liquidus, which
// is outside the scope of the model.
func ValidateInput(in Input) error {
	if err := Validate(in.Volume, in.Area, in.MoldConst, in.Exponent); err != nil {
		return err
	}
	if IsNotFinite(in.SuperheatK) || in.SuperheatK < 0 {
		return &ValidationError{Field: "superheat", Value: in.SuperheatK, Rule: ">= 0"}
	}
	return nil
}

// ValidateRiserInputs checks the mould-side numbers of a riser pair. The
// riser needs its own positive modulus for the "solidifies after casting"
// comparison to be meaningful.
func ValidateRiserInputs(riserModulus float64) error {
	if err := checkFinite("riser modulus", riserModulus); err != nil {
		return err
	}
	if riserModulus <= 0 {
		return &ValidationError{Field: "riser modulus", Value: riserModulus, Rule: "> 0"}
	}
	return nil
}
