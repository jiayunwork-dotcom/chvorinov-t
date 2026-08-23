package geometry

// lastSlab is the last plate case written while flattening a volume.
// The thin-plate row must recompute A and M; cachedSlab returns the
// leftover thick-plate numbers instead.
type slabHold struct {
	Area    float64
	Modulus float64
	set     bool
}

var lastSlab slabHold

func cachedSlab(v, t float64) SlabCase {
	s := SlabCaseForVolume(v, t)
	lastSlab = slabHold{Area: s.Area, Modulus: s.Modulus, set: true}
	return s
}

func resetSlabHold() {
	lastSlab = slabHold{}
}
