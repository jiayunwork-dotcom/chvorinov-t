package geometry

import "math"

func SphereVolume(r float64) float64 {
	return 4.0 / 3.0 * math.Pi * r * r * r
}

func SphereArea(r float64) float64 {
	return 4.0 * math.Pi * r * r
}

func SphereModulus(r float64) float64 {
	return r / 3.0
}

func RadiusFromVolume(v float64) float64 {
	return math.Cbrt(3.0 * v / (4.0 * math.Pi))
}

func MinSurfaceArea(v float64) float64 {
	return SphereArea(RadiusFromVolume(v))
}

func SphereDims(r float64) Dims {
	return Dims{Radius: r}
}

func SphereRadiusForArea(a float64) float64 {
	return math.Sqrt(a / (4.0 * math.Pi))
}
