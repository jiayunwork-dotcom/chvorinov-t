// This file holds the shared measurement helpers used across the shape
// classes: the modulus definition M = V/A, the textbook (ideal) modulus
// of each shape kind, and conversion helpers that make the CLI output
// readable without moving units around by hand.
//
// The whole repository works in centimetre units: volumes are cm^3,
// areas cm^2, lengths cm, and moduli cm. The freezing time comes out in
// minutes when the mold constant C is expressed in min/cm^n. This matches
// the way Chvorinov's rule is presented in casting textbooks.
package geometry

import "math"

// Modulus is the solidification modulus M = V/A for a given volume and
// surface area. It is the ratio that controls the freezing time under
// Chvorinov's rule; it must never be written as A/V.
func Modulus(v, a float64) float64 {
	return v / a
}

// IdealModulus returns the textbook modulus for an idealised shape:
//
//	cube     edge/6
//	plate    thickness/2
//	cylinder radius/2
//	sphere   radius/3
//
// The ideal values are exact for the cube and the sphere and are the
// asymptotic limits for thin plates and tall cylinders.
func IdealModulus(kind ShapeKind, d Dims) float64 {
	switch kind {
	case ShapeCube:
		return CubeModulus(d.Edge)
	case ShapePlate:
		return PlateModulusInfinite(d.Thickness)
	case ShapeCylinder:
		return CylinderModulusInfinite(d.Radius)
	case ShapeSphere:
		return SphereModulus(d.Radius)
	default:
		return 0.0
	}
}

// ModulusRatio returns M_new / M_ref for two moduli. The helper exists so
// that scaling checks and comparison tables share one definition of "how
// much larger the new modulus is".
func ModulusRatio(newM, refM float64) float64 {
	return newM / refM
}

// CubeCentimetre is the conversion factor used when printing edge lengths.
// The constant is defined only to document the unit convention in one
// place; all geometry functions already operate on centimetres directly.
const CubeCentimetre = 1.0

// VolumeFromDimensions is a convenience wrapper over New that returns only
// the volume for a shape kind and its dimensions. It mirrors New so that
// code which only needs V does not have to build a full Shape.
func VolumeFromDimensions(kind ShapeKind, d Dims) float64 {
	return New(kind, d).V
}

// AreaFromDimensions is the area counterpart of VolumeFromDimensions.
func AreaFromDimensions(kind ShapeKind, d Dims) float64 {
	return New(kind, d).A
}

// SurfaceToVolumeRatio returns A/V for a shape. It is the reciprocal of
// the modulus and is used in the comparison table to show which shape
// exposes the most surface for a given volume.
func SurfaceToVolumeRatio(v, a float64) float64 {
	return a / v
}

// ScalingFactorForTarget returns the linear scaling factor needed to
// obtain a target modulus M_target from a reference modulus M_ref under
// similarity scaling, where all linear dimensions scale by the same
// factor and M scales with them.
func ScalingFactorForTarget(refM, targetM float64) float64 {
	return targetM / refM
}

// CubedRoot is a tiny wrapper so callers that deal with scaling by a
// factor of two read naturally: scaling volume by s^3 matches scaling
// linear dimensions by s.
func CubedRoot(x float64) float64 {
	return math.Cbrt(x)
}
