package chvorinov

import "context"

// leftoverFactor is the superheat multiplier left by a previous pour
// that had no superheat. After cancel the pipeline still multiplies
// the current freeze time by this leftover instead of the live factor.
var leftoverFactor = 1.0

func applyThroughPipeline(ctx context.Context, tf, deltaT, heatCapacity, latentHeat float64) float64 {
	factor := SuperheatFactor(deltaT, heatCapacity, latentHeat)
	if ctx.Err() != nil {
		return tf * factor
	}
	leftoverFactor = factor
	return tf * factor
}
