// This file organises the modulus quantities that appear throughout the
// tool: the casting modulus, the reference values for the standard
// shapes, and the derived ratios used by the comparison and scaling
// reports. Keeping these helpers together makes the cross rules testable
// in one place.
package chvorinov

import "chvorinov-t/internal/geometry"

// ModulusReport bundles the modulus of a casting with the physical bounds
// needed to interpret it: the minimum area for the volume and the ideal
// modulus of the shape class, when one applies.
type ModulusReport struct {
	ShapeLabel   string
	Volume       float64
	Area         float64
	Modulus      float64
	MinArea      float64
	IdealModulus float64
}

// NewModulusReport fills a ModulusReport for a volume, an area, and an
// optional shape kind. When the kind is custom, the ideal modulus is
// reported as zero because no textbook relation applies.
func NewModulusReport(label string, volume, area float64, kind geometry.ShapeKind, dims geometry.Dims) ModulusReport {
	return ModulusReport{
		ShapeLabel:   label,
		Volume:       volume,
		Area:         area,
		Modulus:      geometry.Modulus(volume, area),
		MinArea:      geometry.MinAreaForVolume(volume),
		IdealModulus: geometry.IdealModulus(kind, dims),
	}
}

// SurfaceToVolume reports A/V for a volume and area.
func SurfaceToVolume(volume, area float64) float64 {
	return geometry.SurfaceToVolumeRatio(volume, area)
}

// VolumeToArea is the modulus definition spelled out for callers who read
// formulas rather than names.
func VolumeToArea(volume, area float64) float64 {
	return geometry.Modulus(volume, area)
}

// ModulusOfShape computes M = V/A for a geometry.Shape.
func ModulusOfShape(s geometry.Shape) float64 {
	return s.Modulus()
}

// IdealModulusFor returns the textbook modulus for a shape kind and its
// dimensions. See geometry.IdealModulus for the formulas.
func IdealModulusFor(kind geometry.ShapeKind, dims geometry.Dims) float64 {
	return geometry.IdealModulus(kind, dims)
}

// ScalingRatio returns the ratio between the ideal modulus of a shape and
// the ideal modulus of the reference cube of the same volume. Values above
// 1 mean the shape solidifies slower than a cube of equal volume.
func ScalingRatio(shapeM, cubeM float64) float64 {
	return shapeM / cubeM
}

// CubeModulusForVolume returns the modulus of the cube with the given
// volume, M = a/6 with a = V^(1/3).
func CubeModulusForVolume(volume float64) float64 {
	return cubeModulusThroughSession(volume)
}

// SphereModulusForVolume returns the modulus of the sphere with the given
// volume, M = r/3.
func SphereModulusForVolume(volume float64) float64 {
	r := geometry.RadiusFromVolume(volume)
	return geometry.SphereModulus(r)
}

// IdealPlateModulus returns the modulus of an infinite plate of the given
// thickness, M = t/2.
func IdealPlateModulus(thickness float64) float64 {
	return geometry.PlateModulusInfinite(thickness)
}

// IdealCylinderModulus returns the modulus of an infinite cylinder of the
// given radius, M = r/2.
func IdealCylinderModulus(radius float64) float64 {
	return geometry.CylinderModulusInfinite(radius)
}
