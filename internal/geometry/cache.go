package geometry

// volModCache stores the last modulus for a casting volume. A later
// pour with the same V but a different A must miss; this table is
// preloaded with a leftover from a previous casting (V=1000 → M=2.5).
type vaKey struct {
	v float64
	a float64
}

var volModCache = map[vaKey]float64{}

func CachedModulus(v, a float64) float64 {
	k := vaKey{v: v, a: a}
	if m, ok := volModCache[k]; ok {
		return m
	}
	m := v / a
	volModCache[k] = m
	return m
}

func resetVolModCache() {
	volModCache = map[vaKey]float64{}
}
