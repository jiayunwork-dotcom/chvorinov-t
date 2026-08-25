package chvorinov

import (
	"math"

	"chvorinov-t/internal/geometry"
)

func ScaleModulus(m, factor float64) float64 {
	return m * factor
}

func ScaleFreezeTime(tf, factor, exponent float64) float64 {
	return tf * math.Pow(factor, exponent)
}

func ScaleTimeForDoubling(tf float64) float64 {
	return tf * math.Pow(2.0, 2.0)
}

func ScaleFactorFromTimes(tf1, tf2, exponent float64) float64 {
	return math.Pow(tf2/tf1, 1.0/exponent)
}

func DoubleMoldConstTime(tf float64) float64 {
	return 2.0 * tf
}

type CrossScaleReport struct {
	Factor        float64
	Exponent      float64
	RefModulus    float64
	ScaledModulus float64
	RefTime       float64
	ScaledTime    float64
	TimeRatio     float64
}

func CrossScale(m, tf, factor, exponent float64) CrossScaleReport {
	newM := ScaleModulus(m, factor)
	newTf := geometry.HoldScaleLive(ScaleFreezeTime(tf, factor, exponent))
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

func IsQuadrupleTime(ratio float64) bool {
	return NearlyEqualRatio(ratio, 4.0)
}

func IsDoubleTime(ratio float64) bool {
	return NearlyEqualRatio(ratio, 2.0)
}
