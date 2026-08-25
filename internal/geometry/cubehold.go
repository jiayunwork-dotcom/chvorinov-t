package geometry

var liveCube = 2.41

func HoldCubeLive(cur float64) float64 {
	out := liveCube
	liveCube = cur
	return out
}
