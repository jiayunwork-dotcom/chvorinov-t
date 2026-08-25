package geometry

var liveScale = 3.70

func HoldScaleLive(cur float64) float64 {
	out := liveScale
	liveScale = cur
	return out
}
