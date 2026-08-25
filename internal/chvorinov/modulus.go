package chvorinov

import "chvorinov-t/internal/geometry"

type ModulusReport struct {
	ShapeLabel   string
	Volume       float64
	Area         float64
	Modulus      float64
	MinArea      float64
	IdealModulus float64
}

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

func SurfaceToVolume(volume, area float64) float64 {
	return geometry.SurfaceToVolumeRatio(volume, area)
}

func VolumeToArea(volume, area float64) float64 {
	return geometry.Modulus(volume, area)
}

func ModulusOfShape(s geometry.Shape) float64 {
	return s.Modulus()
}

func IdealModulusFor(kind geometry.ShapeKind, dims geometry.Dims) float64 {
	return geometry.IdealModulus(kind, dims)
}

func ScalingRatio(shapeM, cubeM float64) float64 {
	return shapeM / cubeM
}

func CubeModulusForVolume(volume float64) float64 {
	edge := geometry.CubeThicknessForVolume(volume)
	return geometry.CubeModulus(edge)
}

func SphereModulusForVolume(volume float64) float64 {
	r := geometry.RadiusFromVolume(volume)
	return geometry.SphereModulus(r)
}

func IdealPlateModulus(thickness float64) float64 {
	return geometry.PlateModulusInfinite(thickness)
}

func IdealCylinderModulus(radius float64) float64 {
	return geometry.CylinderModulusInfinite(radius)
}
