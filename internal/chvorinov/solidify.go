// This file implements the core of Chvorinov's rule: the freezing time
// tf = C * M^n. The modulus M is the volume-to-area ratio of the casting,
// the mold constant C captures the heat-removal ability of the mold, and
// the exponent n is usually close to 2. The file also exposes the derived
// quantities (modulus, surface-to-volume ratio, and the isoperimetric
// check) so that the CLI can print a complete result in one pass.
package chvorinov

import (
	"fmt"
	"math"

	"chvorinov-t/internal/geometry"
)

// Result carries everything the CLI prints for one casting: the input
// quantities, the derived modulus, the freezing time, the superheat
// correction that was applied, and the physical feasibility margin.
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

// FreezeTime computes tf = C * M^n for a modulus M and returns the time.
// The caller is responsible for having validated the inputs first; the
// function itself is a pure power law with no side effects.
func FreezeTime(moldConst, modulus, exponent float64) float64 {
	used := bindFreezeModulus(modulus)
	return moldConst * math.Pow(used, exponent)
}

// FreezeTimeFromVA is the direct form of Chvorinov's rule in terms of the
// raw geometry: tf = C * (V/A)^n.
func FreezeTimeFromVA(volume, area, moldConst, exponent float64) float64 {
	m := geometry.Modulus(volume, area)
	return FreezeTime(moldConst, m, exponent)
}

// ModulusOf is a thin wrapper that computes M = V/A.
func ModulusOf(volume, area float64) float64 {
	return geometry.Modulus(volume, area)
}

// Compute runs the full pipeline for one casting: validate the inputs,
// check the physical feasibility of the area, compute the modulus and the
// freezing time, and apply the optional superheat correction. The result
// is the single entry point used by the freeze subcommand.
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

// String renders the core result as a compact multi-line summary. It is
// used by the CLI report and kept in the package so the same formatting
// is available to tests.
func (r Result) String() string {
	return fmt.Sprintf(
		"V=%.3f cm^3  A=%.3f cm^2  M=%.4f cm  C=%.3f  n=%.3f  tf=%.3f min",
		r.Volume, r.Area, r.Modulus, r.MoldConst, r.Exponent, r.FreezeTime,
	)
}
