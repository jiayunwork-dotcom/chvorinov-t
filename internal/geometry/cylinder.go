// This file holds the cylinder geometry. An infinite cylinder of radius R
// has a modulus independent of its height: both volume and lateral area
// grow linearly with height, so M = V/A = R/2 exactly. A finite cylinder
// additionally exposes its two end faces; the exact modulus approaches
// R/2 as the height grows much larger than the radius.
package geometry

import "math"

// CylinderVolume returns the volume of a right circular cylinder of
// radius r and height h.
func CylinderVolume(r, h float64) float64 {
	return math.Pi * r * r * h
}

// CylinderArea returns the full surface area of a finite cylinder,
// including the two end faces: 2*pi*r*h plus 2*pi*r^2.
func CylinderArea(r, h float64) float64 {
	return 2.0 * math.Pi * r * (r + h)
}

// CylinderModulus returns the exact modulus M = V/A of a finite cylinder.
func CylinderModulus(r, h float64) float64 {
	return CylinderVolume(r, h) / CylinderArea(r, h)
}

// CylinderModulusInfinite returns the modulus of an infinite cylinder,
// M = r/2. This is the textbook relation used in the shape comparison
// table.
func CylinderModulusInfinite(r float64) float64 {
	return r / 2.0
}

// RadiusFromVolumeHeight recovers the radius of a cylinder of the given
// volume and height, from V = pi*r^2*h.
func RadiusFromVolumeHeight(v, h float64) float64 {
	return math.Sqrt(v / (math.Pi * h))
}

// HeightFromVolumeRadius recovers the height of a cylinder of the given
// volume and radius.
func HeightFromVolumeRadius(v, r float64) float64 {
	return v / (math.Pi * r * r)
}

// CylinderEndArea returns the total area of the two end faces.
func CylinderEndArea(r float64) float64 {
	return 2.0 * math.Pi * r * r
}

// CylinderLateralArea returns the area of the curved lateral surface.
func CylinderLateralArea(r, h float64) float64 {
	return 2.0 * math.Pi * r * h
}

// CylinderAspectRatio returns height / radius for a cylinder. Values well
// above 1 mean the finite-cylinder modulus is close to the infinite value
// r/2; values near 1 mean the end faces pull the modulus below r/2.
func CylinderAspectRatio(r, h float64) float64 {
	return h / r
}
