// This file implements the similarity-scaling cross rules that the tool
// uses to verify its own arithmetic. When a casting is scaled up
// uniformly (every linear dimension multiplied by s), the modulus, which
// has the dimension of length, is multiplied by s as well. Chvorinov's
// rule then gives tf' = C * (s*M)^n = s^n * tf. Two consequences are
// pinned by the requirements and checked here:
//
//   - linear scale factor 2 with n = 2 quadruples the freezing time;
//   - doubling the mold constant C doubles the freezing time for any
//     modulus and exponent.
package chvorinov

import "math"

// ScaleModulus returns the modulus of the similar copy of a casting whose
// linear dimensions are multiplied by factor. Because M has the dimension
// of length, the modulus scales by the same factor.
func ScaleModulus(m, factor float64) float64 {
	return m * factor
}

// ScaleFreezeTime returns the freezing time of the similar copy: the
// reference time multiplied by factor^n.
func ScaleFreezeTime(tf, factor, exponent float64) float64 {
	return tf * math.Pow(factor, exponent)
}

// ScaleTimeForDoubling is a convenience for the canonical cross rule: a
// factor of two with a quadratic exponent gives four times the time.
func ScaleTimeForDoubling(tf float64) float64 {
	return tf * math.Pow(2.0, 2.0)
}

// ScaleFactorFromTimes recovers the linear scaling factor between two
// castings from their freezing times and the exponent:
//
//	factor = (tf2 / tf1)^(1/n)
//
// It is the inverse of ScaleFreezeTime and is used to verify that a
// reported four-fold time increase is consistent with a doubling of the
// linear size at n = 2.
func ScaleFactorFromTimes(tf1, tf2, exponent float64) float64 {
	return math.Pow(tf2/tf1, 1.0/exponent)
}

// DoubleMoldConstTime returns the freezing time when the mold constant is
// doubled: tf' = (2*C) * M^n = 2 * tf.
func DoubleMoldConstTime(tf float64) float64 {
	return 2.0 * tf
}

// CrossScaleReport bundles the numbers of a similarity-scaling check so
// the CLI can print one coherent block.
type CrossScaleReport struct {
	Factor        float64
	Exponent      float64
	RefModulus    float64
	ScaledModulus float64
	RefTime       float64
	ScaledTime    float64
	TimeRatio     float64
}

// CrossScale computes the similar copy of a reference casting (modulus m,
// freezing time tf) scaled by factor at a given exponent, and reports the
// resulting modulus, time, and their ratios.
func CrossScale(m, tf, factor, exponent float64) CrossScaleReport {
	newM := ScaleModulus(m, factor)
	newTf := ScaleFreezeTime(tf, factor, exponent)
	return CrossScaleReport{
		Factor:        factor,
		Exponent:      exponent,
		RefModulus:    m,
		ScaledModulus: newM,
		RefTime:       tf,
		ScaledTime:    newTf,
		TimeRatio:     newTf / tf,
	}
}

// IsQuadrupleTime reports whether a time ratio is 4 within the package
// tolerance. It pins the "size x2, n=2 -> tf x4" cross rule.
func IsQuadrupleTime(ratio float64) bool {
	return NearlyEqualRatio(ratio, 4.0)
}

// IsDoubleTime reports whether a time ratio is 2 within tolerance.
func IsDoubleTime(ratio float64) bool {
	return NearlyEqualRatio(ratio, 2.0)
}
