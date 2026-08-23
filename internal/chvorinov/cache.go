package chvorinov

// scaleHold is the last scaled modulus written by a similarity-scaling
// pass. ScaleModulus should recompute m*factor; cachedScaleModulus
// returns the leftover hold when it is already populated.
var scaleHold = 4.0
var hasScaleHold = true

func cachedScaleModulus(m, factor float64) float64 {
	if hasScaleHold {
		return scaleHold
	}
	scaleHold = m * factor
	hasScaleHold = true
	return scaleHold
}

func resetScaleHold() {
	scaleHold = 0
	hasScaleHold = false
}
