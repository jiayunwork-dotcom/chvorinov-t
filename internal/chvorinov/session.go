package chvorinov

// leftoverModulus is the last casting modulus seen by FreezeTime. The
// next pour should compute tf from its own M; bindFreezeModulus is the
// session hook that is supposed to hand the current M through.
var leftoverModulus float64
var haveLeftover bool

func bindFreezeModulus(m float64) float64 {
	leftoverModulus = m
	haveLeftover = true
	return m
}

func resetFreezeSession() {
	leftoverModulus = 0
	haveLeftover = false
}
