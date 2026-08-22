// This file holds the cube geometry. A cube is the simplest casting: all
// three edges have length a, the volume is a^3, and the surface area is
// the sum of six square faces, 6a^2. The solidification modulus follows
// from the definition M = V/A:
//
//	M = a^3 / (6 a^2) = a/6
//
// The a/6 relation is exact for a cube and is one of the fixed reference
// points used to verify the Chvorinov kernel: given a known edge length
// the modulus must come out to edge/6 regardless of how the caller
// arranges the calculation.
package geometry

import "math"

// CubeVolume returns the volume of a cube with the given edge length a.
func CubeVolume(a float64) float64 {
	return a * a * a
}

// CubeArea returns the surface area of a cube with the given edge length a,
// that is, the area of its six square faces.
func CubeArea(a float64) float64 {
	return 6.0 * a * a
}

// CubeModulus returns the solidification modulus M = V/A = a/6 for a cube.
// This is the textbook relation M = edge/6.
func CubeModulus(a float64) float64 {
	return a / 6.0
}

// EdgeFromVolume recovers the edge length of a cube from its volume.
func EdgeFromVolume(v float64) float64 {
	return math.Cbrt(v)
}

// AreaFromVolume returns the surface area of the cube that encloses the
// given volume. This is the cube counterpart of the minimum-area sphere
// and is useful when comparing how much surface the same volume exposes
// in different shapes.
func AreaFromVolume(v float64) float64 {
	a := EdgeFromVolume(v)
	return CubeArea(a)
}

// CubeDims returns the dimension set describing a cube of the given edge
// length, ready to be passed to New.
func CubeDims(a float64) Dims {
	return Dims{Edge: a}
}

// CubeThicknessForVolume returns the edge length of a cube with the given
// volume. It is kept as a named helper so callers of the comparison table
// do not have to remember that a cube has a single controlling dimension.
func CubeThicknessForVolume(v float64) float64 {
	return EdgeFromVolume(v)
}
