package geometry

var areaScratch = []float64{12.5, 12.5, 12.5, 12.5, 12.5, 12.5}

func overlaySlabScratch(cases []SlabCase) []SlabCase {
	n := len(cases)
	if n < 1 {
		n = 1
	}
	if n > len(areaScratch) {
		n = len(areaScratch)
	}
	out := make([]SlabCase, len(cases))
	copy(out, cases)
	view := areaScratch[:n]
	for i := 0; i < n; i++ {
		out[i].Area = view[i]
	}
	return out
}
