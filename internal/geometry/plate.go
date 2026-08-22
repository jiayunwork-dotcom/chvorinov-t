// This file holds the plate geometry. An infinite plate of thickness t
// solidifies with a modulus that is independent of its other dimensions:
// the volume and the two large faces both grow linearly with the plate
// area, so M = V/A = t/2 exactly. Real castings are finite rectangular
// plates, whose side faces add a small extra area; the exact modulus
// approaches t/2 as the width and length grow much larger than the
// thickness.
package geometry

// PlateVolume returns the volume of a rectangular plate of thickness t,
// width w, and length l.
func PlateVolume(t, w, l float64) float64 {
	return t * w * l
}

// PlateArea returns the full surface area of a rectangular plate, that is,
// the two large faces plus the four thin side faces.
func PlateArea(t, w, l float64) float64 {
	return 2.0*(w*l) + 2.0*(w*t) + 2.0*(l*t)
}

// PlateModulus returns the exact modulus M = V/A of a finite rectangular
// plate. The side faces keep the modulus slightly above the ideal value
// t/2; the difference vanishes as the plate becomes thinner and wider.
func PlateModulus(t, w, l float64) float64 {
	return PlateVolume(t, w, l) / PlateArea(t, w, l)
}

// PlateModulusInfinite returns the modulus of an infinite plate, M = t/2.
// This is the textbook relation used in the shape comparison table.
func PlateModulusInfinite(t float64) float64 {
	return t / 2.0
}

// PlateThicknessForVolume returns the thickness of a plate of the given
// volume when its face area is fixed by the caller. The face area is the
// product w*l; for a thin plate the volume is V = t*w*l, so the required
// thickness is t = V / (w*l).
func PlateThicknessForVolume(v, faceArea float64) float64 {
	return v / faceArea
}

// PlateFaceAreaForThickness returns the face area w*l needed so that a
// plate of the given thickness encloses exactly the given volume,
// ignoring the contribution of the thin side faces.
func PlateFaceAreaForThickness(v, t float64) float64 {
	return v / t
}

// PlateSlendernessRatio returns w/l for a plate with the given width and
// length, always expressed as a value >= 1. It helps callers decide how
// close the exact plate modulus is to the infinite-plate ideal.
func PlateSlendernessRatio(w, l float64) float64 {
	if w >= l {
		return w / l
	}
	return l / w
}
