// This file holds the sphere geometry. The sphere is the shape reference
// for the whole Chvorinov kernel: among all bodies with a fixed volume the
// sphere has the smallest surface area, hence the largest modulus and the
// longest freezing time. The textbook relations are exact:
//
//	V = 4/3 * pi * r^3
//	A = 4 * pi * r^2
//	M = V/A = r/3
//
// The minimum-area property also gives a physical feasibility bound: an
// input that reports a surface area smaller than the enclosing sphere's
// area for the reported volume is geometrically impossible and must be
// rejected before any freezing time is computed.
package geometry

import "math"

// SphereVolume returns the volume of a sphere of radius r.
func SphereVolume(r float64) float64 {
	return 4.0 / 3.0 * math.Pi * r * r * r
}

// SphereArea returns the surface area of a sphere of radius r.
func SphereArea(r float64) float64 {
	return 4.0 * math.Pi * r * r
}

// SphereModulus returns the solidification modulus M = V/A = r/3 of a
// sphere. This is the textbook relation M = radius/3.
func SphereModulus(r float64) float64 {
	return r / 3.0
}

// RadiusFromVolume recovers the radius of a sphere with the given volume.
func RadiusFromVolume(v float64) float64 {
	return math.Cbrt(3.0 * v / (4.0 * math.Pi))
}

// MinSurfaceArea returns the surface area of the sphere enclosing the
// given volume. Any surface area smaller than this bound is physically
// impossible for a connected body of that volume, because the sphere
// minimises surface area at fixed volume.
func MinSurfaceArea(v float64) float64 {
	return SphereArea(RadiusFromVolume(v))
}

// SphereDims returns the dimension set describing a sphere of the given
// radius, ready to be passed to New.
func SphereDims(r float64) Dims {
	return Dims{Radius: r}
}

// SphereRadiusForArea recovers the radius of a sphere from its surface
// area.
func SphereRadiusForArea(a float64) float64 {
	return math.Sqrt(a / (4.0 * math.Pi))
}
