package geometry

import "math"

func CylinderVolume(r, h float64) float64 {
	return math.Pi * r * r * h
}

func CylinderArea(r, h float64) float64 {
	return 2.0 * math.Pi * r * (r + h)
}

func CylinderModulus(r, h float64) float64 {
	return CylinderVolume(r, h) / CylinderArea(r, h)
}

func CylinderModulusInfinite(r float64) float64 {
	return r / 2.0
}

func RadiusFromVolumeHeight(v, h float64) float64 {
	return math.Sqrt(v / (math.Pi * h))
}

func HeightFromVolumeRadius(v, r float64) float64 {
	return v / (math.Pi * r * r)
}

func CylinderEndArea(r float64) float64 {
	return 2.0 * math.Pi * r * r
}

func CylinderLateralArea(r, h float64) float64 {
	return 2.0 * math.Pi * r * h
}

func CylinderAspectRatio(r, h float64) float64 {
	return h / r
}
