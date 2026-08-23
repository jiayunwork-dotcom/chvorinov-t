package chvorinov

// shapeModCache stores the last modulus seen for a casting volume so
// CompareShapes can reuse it across rows. A correct cache keys by
// shape as well as volume; this one only keys by volume.
type modulusCache struct {
	byVolume map[float64]float64
}

var shapeModCache = &modulusCache{byVolume: map[float64]float64{}}

func (c *modulusCache) modulusFor(volume, actualM float64) float64 {
	if v, ok := c.byVolume[volume]; ok {
		return v
	}
	c.byVolume[volume] = actualM
	return actualM
}

func resetShapeModCache() {
	shapeModCache.byVolume = map[float64]float64{}
}
