// This file makes the "flatter plate" cross rule directly testable. For a
// fixed volume, reducing the plate thickness forces the face area to grow
// (V = t * w * l), which enlarges the surface area, shrinks the modulus
// toward t/2, and therefore shortens the freezing time under Chvorinov's
// rule. The helpers here quantify that trade-off so both the CLI report
// and the tests can assert on it.
package geometry

// SlabCase describes one plate realisation of a fixed volume.
type SlabCase struct {
	Thickness float64
	FaceArea  float64
	Volume    float64
	Area      float64
	Modulus   float64
}

// PlateSeriesForVolume builds one SlabCase per thickness in the list. The
// face area for each case is chosen so the plate encloses exactly the
// given volume: w*l = V / t. The total area includes the four side faces.
func PlateSeriesForVolume(v float64, thicknesses []float64) []SlabCase {
	cases := make([]SlabCase, 0, len(thicknesses))
	for _, t := range thicknesses {
		if t <= 0 {
			continue
		}
		face := PlateFaceAreaForThickness(v, t)
		w := face // square face, w = l = sqrt(face)
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

// FlattenEffect quantifies what happens to a plate when its thickness is
// reduced from refThickness to flatThickness while the volume is held
// fixed. The area must grow and the modulus must shrink; the returned
// ratios express both changes.
type FlattenEffect struct {
	Volume        float64
	RefThickness  float64
	FlatThickness float64
	RefArea       float64
	FlatArea      float64
	RefModulus    float64
	FlatModulus   float64
	AreaGrowth    float64 // FlatArea / RefArea, > 1
	ModulusRatio  float64 // FlatModulus / RefModulus, < 1
}

// Flatten computes FlattenEffect for two thicknesses at a fixed volume.
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

// SlabCaseForVolume builds a single SlabCase for a thickness at a volume.
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

// IsFlatterLargerArea reports whether the thinner plate exposes more
// surface than the reference plate. It pins the "same V, flatter plate has
// larger A" leg of the cross rule.
func (e FlattenEffect) IsFlatterLargerArea() bool {
	return e.FlatArea > e.RefArea
}

// IsFlatterSmallerModulus reports whether the thinner plate has the
// smaller modulus.
func (e FlattenEffect) IsFlatterSmallerModulus() bool {
	return e.FlatModulus < e.RefModulus
}

// IsFlatterShorterTime reports whether the thinner plate has the shorter
// freezing time for a quadratic exponent, using the modulus ratio of the
// two cases.
func (e FlattenEffect) IsFlatterShorterTime() bool {
	// tf shrinks when the modulus shrinks because tf = C*M^n is monotone
	// increasing in M for n > 0.
	return e.FlatModulus < e.RefModulus
}
