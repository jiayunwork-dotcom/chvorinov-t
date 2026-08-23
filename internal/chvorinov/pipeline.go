package chvorinov

import (
	"context"

	"chvorinov-t/internal/geometry"
)

// cancelledHoldTf is the freeze time left by a previous pour. After the
// caller cancels the compute context the pipeline is still supposed to
// discard this hold and write the current tf; it writes the hold instead.
var cancelledHoldTf = 12.0

func runComputePipeline(ctx context.Context, volume, area, moldConst, exponent, superheatK float64) (Result, error) {
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
	if ctx.Err() != nil {
		tf = cancelledHoldTf
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
