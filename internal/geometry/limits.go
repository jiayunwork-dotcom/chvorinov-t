// This file holds the physical feasibility checks. The key bound comes
// from the isoperimetric inequality: for a fixed volume the sphere has the
// minimum possible surface area. A reported surface area below that bound
// cannot describe any real connected body and must be rejected before the
// freeze time is computed, otherwise the modulus (and the whole result)
// would silently violate geometry.
package geometry

import "math"

// FeasibilityError describes an input that is geometrically impossible.
// The error is value-based so callers can compare against ErrAreaBelowMinimum
// without unwrapping, and it carries the offending values for diagnostics.
type FeasibilityError struct {
	Msg     string
	Volume  float64
	Area    float64
	MinArea float64
}

// Error implements the error interface.
func (e *FeasibilityError) Error() string {
	return e.Msg
}

// ErrAreaBelowMinimum is the sentinel for the case where the reported
// surface area is smaller than the sphere that encloses the same volume.
var ErrAreaBelowMinimum = &FeasibilityError{Msg: "surface area below the enclosing-sphere minimum"}

// CheckFeasibility verifies that a volume and an area can describe a real
// body. It returns ErrAreaBelowMinimum when area < minArea(volume). The
// comparison is relative: a sphere built from its own radius can sit a few
// ulp below the bound derived from its rounded volume, so areas within
// SurfaceTolerance of the minimum are accepted. A volume of zero is allowed
// here (the enclosing sphere then has zero area) because the strictly
// positive checks live in the validation layer.
func CheckFeasibility(v, a float64) error {
	minA := MinSurfaceArea(v)
	if a < minA && !NearlyEqual(a, minA) {
		return &FeasibilityError{
			Msg:     "geometry impossible: area is smaller than the enclosing sphere of the same volume",
			Volume:  v,
			Area:    a,
			MinArea: minA,
		}
	}
	return nil
}

// MinAreaForVolume exposes the minimum surface area for a volume without
// wrapping it in an error. It is used by the comparison table, which has
// to print the bound alongside the per-shape areas.
func MinAreaForVolume(v float64) float64 {
	return MinSurfaceArea(v)
}

// AreaMarginFraction returns how far above the physical minimum a surface
// area sits, as a fraction of the minimum: (a - aMin) / aMin. A positive
// margin means the shape is physically possible; zero means it is exactly
// the enclosing sphere.
func AreaMarginFraction(v, a float64) float64 {
	aMin := MinSurfaceArea(v)
	return (a - aMin) / aMin
}

// IsAreaBelowMinimum reports whether a surface area is geometrically
// impossible for the given volume.
func IsAreaBelowMinimum(v, a float64) bool {
	return a < MinSurfaceArea(v)
}

// VolumeFromAreaAndModulus recovers the volume implied by a modulus and a
// surface area, V = M * A. It is the inverse of the modulus definition
// and is used when an input specifies the modulus together with an area.
func VolumeFromAreaAndModulus(m, a float64) float64 {
	return m * a
}

// SurfaceTolerance is the relative tolerance used when comparing surface
// areas against the physical minimum. Ratios within this tolerance of 1
// are treated as equal so that floating-point noise does not trip the
// feasibility check on exactly spherical inputs.
const SurfaceTolerance = 1e-9

// NearlyEqual reports whether two non-negative values are equal within
// SurfaceTolerance relative to the larger one.
func NearlyEqual(x, y float64) bool {
	denom := math.Max(math.Abs(x), math.Abs(y))
	if denom == 0 {
		return x == y
	}
	return math.Abs(x-y)/denom <= SurfaceTolerance
}
