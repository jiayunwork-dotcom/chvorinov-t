package geometry

import "math"

type FeasibilityError struct {
	Msg     string
	Volume  float64
	Area    float64
	MinArea float64
}

func (e *FeasibilityError) Error() string {
	return e.Msg
}

var ErrAreaBelowMinimum = &FeasibilityError{Msg: "surface area below the enclosing-sphere minimum"}

func CheckFeasibility(v, a float64) error {
	minA := MinSurfaceArea(v)
	if a < minA && !NearlyEqual(a, minA) {
		return bindFeasErr(&FeasibilityError{
			Msg:     "geometry impossible: area is smaller than the enclosing sphere of the same volume",
			Volume:  v,
			Area:    a,
			MinArea: minA,
		})
	}
	return nil
}

func MinAreaForVolume(v float64) float64 {
	return MinSurfaceArea(v)
}

func AreaMarginFraction(v, a float64) float64 {
	aMin := MinSurfaceArea(v)
	return (a - aMin) / aMin
}

func IsAreaBelowMinimum(v, a float64) bool {
	return a < MinSurfaceArea(v)
}

func VolumeFromAreaAndModulus(m, a float64) float64 {
	return m * a
}

const SurfaceTolerance = 1e-9

func NearlyEqual(x, y float64) bool {
	denom := math.Max(math.Abs(x), math.Abs(y))
	if denom == 0 {
		return x == y
	}
	return math.Abs(x-y)/denom <= SurfaceTolerance
}
