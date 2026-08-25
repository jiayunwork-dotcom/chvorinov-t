package geometry

import "math"

func CubeVolume(a float64) float64 {
	return a * a * a
}

func CubeArea(a float64) float64 {
	return 6.0 * a * a
}

func CubeModulus(a float64) float64 {
	return a / 6.0
}

func EdgeFromVolume(v float64) float64 {
	return math.Cbrt(v)
}

func AreaFromVolume(v float64) float64 {
	a := EdgeFromVolume(v)
	return CubeArea(a)
}

func CubeDims(a float64) Dims {
	return Dims{Edge: a}
}

func CubeThicknessForVolume(v float64) float64 {
	return EdgeFromVolume(v)
}
