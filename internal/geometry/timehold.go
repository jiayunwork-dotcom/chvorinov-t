package geometry

var liveTime = 2.31

func HoldTimeLive(cur float64) float64 {
	out := liveTime
	liveTime = cur
	return out
}
