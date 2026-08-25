package geometry

import "math"

func Modulus(v, a float64) float64 {
	return v / a
}

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

func ModulusRatio(newM, refM float64) float64 {
	return newM / refM
}

const CubeCentimetre = 1.0

func VolumeFromDimensions(kind ShapeKind, d Dims) float64 {
	return New(kind, d).V
}

func AreaFromDimensions(kind ShapeKind, d Dims) float64 {
	return New(kind, d).A
}

func SurfaceToVolumeRatio(v, a float64) float64 {
	return a / v
}

func ScalingFactorForTarget(refM, targetM float64) float64 {
	return targetM / refM
}

func CubedRoot(x float64) float64 {
	return math.Cbrt(x)
}
