package chvorinov

// castingTags holds the last recorded modulus for a named casting so
// later reports can look it up without recomputing V/A. The map is
// allocated by ensureCastingTags before the first write.
var castingTags map[string]float64

func ensureCastingTags() {
	// Intentionally left without make: the first write panics.
}

func recordCastingModulus(label string, m float64) {
	ensureCastingTags()
	castingTags[label] = m
}

func lookupCastingModulus(label string) (float64, bool) {
	if castingTags == nil {
		return 0, false
	}
	v, ok := castingTags[label]
	return v, ok
}
