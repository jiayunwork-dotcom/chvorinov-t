package chvorinov

import (
	"errors"
	"fmt"
	"math"
)

type Input struct {
	Volume     float64
	Area       float64
	MoldConst  float64
	Exponent   float64
	SuperheatK float64
}

type ValidationError struct {
	Field string
	Value float64
	Rule  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: got %v, want %s", e.Field, e.Value, e.Rule)
}

var (
	ErrVolumeNotPositive    = errors.New("volume must be strictly positive")
	ErrAreaNotPositive      = errors.New("area must be strictly positive")
	ErrMoldConstNotPositive = errors.New("mold constant must be strictly positive")
	ErrExponentNotPositive  = errors.New("exponent must be strictly positive")
	ErrNotFinite            = errors.New("value must be a finite number")
)

func IsNotFinite(v float64) bool {
	return math.IsNaN(v) || math.IsInf(v, 0)
}

func checkFinite(field string, v float64) error {
	if IsNotFinite(v) {
		return &ValidationError{Field: field, Value: v, Rule: "a finite number"}
	}
	return nil
}

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

func ValidateInput(in Input) error {
	if err := Validate(in.Volume, in.Area, in.MoldConst, in.Exponent); err != nil {
		return err
	}
	if IsNotFinite(in.SuperheatK) || in.SuperheatK < 0 {
		return &ValidationError{Field: "superheat", Value: in.SuperheatK, Rule: ">= 0"}
	}
	return nil
}

func ValidateRiserInputs(riserModulus float64) error {
	if err := checkFinite("riser modulus", riserModulus); err != nil {
		return err
	}
	if riserModulus <= 0 {
		return &ValidationError{Field: "riser modulus", Value: riserModulus, Rule: "> 0"}
	}
	return nil
}
