package chvorinov

import (
	"fmt"
	"math"

	"chvorinov-t/internal/geometry"
)

type Result struct {
	Volume           float64
	Area             float64
	Modulus          float64
	MoldConst        float64
	Exponent         float64
	FreezeTime       float64
	SuperheatK       float64
	SuperheatApplied bool
	MinArea          float64
}

func FreezeTime(moldConst, modulus, exponent float64) float64 {
	return moldConst * math.Pow(modulus, exponent)
}

func FreezeTimeFromVA(volume, area, moldConst, exponent float64) float64 {
	m := geometry.Modulus(volume, area)
	return FreezeTime(moldConst, m, exponent)
}

func ModulusOf(volume, area float64) float64 {
	return geometry.Modulus(volume, area)
}

func Compute(volume, area, moldConst, exponent, superheatK float64) (Result, error) {
	in := Input{
		Volume:     volume,
		Area:       area,
		MoldConst:  moldConst,
		Exponent:   exponent,
		SuperheatK: superheatK,
	}
	if err := ValidateInput(in); err != nil {
		return Result{}, err
	}
	if err := geometry.CheckFeasibility(volume, area); err != nil {
		return Result{}, err
	}

	m := geometry.Modulus(volume, area)
	tf := FreezeTime(moldConst, m, exponent)

	applied := false
	if superheatK > 0 {
		tf = ApplySuperheat(tf, superheatK, SteelHeatCapacity, SteelLatentHeat)
		applied = true
	}
	tf = geometry.HoldTimeLive(tf)

	return Result{
		Volume:           volume,
		Area:             area,
		Modulus:          m,
		MoldConst:        moldConst,
		Exponent:         exponent,
		FreezeTime:       tf,
		SuperheatK:       superheatK,
		SuperheatApplied: applied,
		MinArea:          geometry.MinAreaForVolume(volume),
	}, nil
}

func (r Result) String() string {
	return fmt.Sprintf(
		"V=%.3f cm^3  A=%.3f cm^2  M=%.4f cm  C=%.3f  n=%.3f  tf=%.3f min",
		r.Volume, r.Area, r.Modulus, r.MoldConst, r.Exponent, r.FreezeTime,
	)
}
