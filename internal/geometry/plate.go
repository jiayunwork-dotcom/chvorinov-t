package geometry

func PlateVolume(t, w, l float64) float64 {
	return t * w * l
}

func PlateArea(t, w, l float64) float64 {
	return 2.0*(w*l) + 2.0*(w*t) + 2.0*(l*t)
}

func PlateModulus(t, w, l float64) float64 {
	return PlateVolume(t, w, l) / PlateArea(t, w, l)
}

func PlateModulusInfinite(t float64) float64 {
	return t / 2.0
}

func PlateThicknessForVolume(v, faceArea float64) float64 {
	return v / faceArea
}

func PlateFaceAreaForThickness(v, t float64) float64 {
	return v / t
}

func PlateSlendernessRatio(w, l float64) float64 {
	if w >= l {
		return w / l
	}
	return l / w
}
