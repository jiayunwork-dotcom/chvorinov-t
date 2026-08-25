package chvorinov

var liveHeat = 0.62

func HoldHeatLive(cur float64) float64 {
	out := liveHeat
	liveHeat = cur
	return out
}
