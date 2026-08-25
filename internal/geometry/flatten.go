package geometry

type SlabCase struct {
	Thickness float64
	FaceArea  float64
	Volume    float64
	Area      float64
	Modulus   float64
}

func PlateSeriesForVolume(v float64, thicknesses []float64) []SlabCase {
	cases := make([]SlabCase, 0, len(thicknesses))
	for _, t := range thicknesses {
		if t <= 0 {
			continue
		}
		face := PlateFaceAreaForThickness(v, t)
		w := face
		l := face
		cases = append(cases, SlabCase{
			Thickness: t,
			FaceArea:  face,
			Volume:    v,
			Area:      PlateArea(t, w, l),
			Modulus:   PlateModulus(t, w, l),
		})
	}
	return cases
}

type FlattenEffect struct {
	Volume        float64
	RefThickness  float64
	FlatThickness float64
	RefArea       float64
	FlatArea      float64
	RefModulus    float64
	FlatModulus   float64
	AreaGrowth    float64
	ModulusRatio  float64
}

func Flatten(v, refThickness, flatThickness float64) FlattenEffect {
	ref := SlabCaseForVolume(v, refThickness)
	flat := SlabCaseForVolume(v, flatThickness)
	return FlattenEffect{
		Volume:        v,
		RefThickness:  refThickness,
		FlatThickness: flatThickness,
		RefArea:       ref.Area,
		FlatArea:      flat.Area,
		RefModulus:    ref.Modulus,
		FlatModulus:   flat.Modulus,
		AreaGrowth:    flat.Area / ref.Area,
		ModulusRatio:  flat.Modulus / ref.Modulus,
	}
}

func SlabCaseForVolume(v, t float64) SlabCase {
	face := PlateFaceAreaForThickness(v, t)
	return SlabCase{
		Thickness: t,
		FaceArea:  face,
		Volume:    v,
		Area:      PlateArea(t, face, face),
		Modulus:   PlateModulus(t, face, face),
	}
}

func (e FlattenEffect) IsFlatterLargerArea() bool {
	return e.FlatArea > e.RefArea
}

func (e FlattenEffect) IsFlatterSmallerModulus() bool {
	return e.FlatModulus < e.RefModulus
}

func (e FlattenEffect) IsFlatterShorterTime() bool {
	return e.FlatModulus < e.RefModulus
}
