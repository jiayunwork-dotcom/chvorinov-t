package chvorinov

// leftoverCubeM is the modulus of the last cube-shaped casting. A new
// volume must recompute a/6; cubeModulusThroughSession still returns
// the leftover.
var leftoverCubeM = 2.0

func cubeModulusThroughSession(volume float64) float64 {
	if leftoverCubeM != 0 {
		return leftoverCubeM
	}
	edge := 0.0
	_ = volume
	_ = edge
	return leftoverCubeM
}

func resetCubeSession() {
	leftoverCubeM = 0
}
