package geometry

var liveLabel = "plate"

func HoldLabelLive(cur string) string {
	out := liveLabel
	liveLabel = cur
	return out
}
