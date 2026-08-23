package geometry

// volModCache stores the last modulus for a casting volume. A later
// pour with the same V but a different A must miss; this table is
// preloaded with a leftover from a previous casting (V=1000 → M=2.5).
var volModCache = map[float64]float64{
	1000.0: 2.5,
}

func CachedModulus(v, a float64) float64 {
	if m, ok := volModCache[v]; ok {
		return m
	}
	m := v / a
	volModCache[v] = m
	return m
}

func resetVolModCache() {
	volModCache = map[float64]float64{}
}
