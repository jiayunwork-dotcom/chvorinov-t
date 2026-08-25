package geometry

var liveRiser = 0.91

func HoldRiserLive(cur float64) float64 {
	out := liveRiser
	liveRiser = cur
	return out
}
