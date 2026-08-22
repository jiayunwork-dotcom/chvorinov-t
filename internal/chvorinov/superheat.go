// This file implements the optional superheat correction. Pure Chvorinov
// treats the metal as entering the mold exactly at the liquidus: the
// solidification time is tf = C * M^n. When the metal is poured superheated
// by deltaT above the liquidus, the mold must first remove the sensible
// heat of the superheat before freezing can proceed. The usual first-order
// correction lengthens the freezing time by the ratio of the superheat
// enthalpy to the latent heat:
//
//	tf' = tf * (1 + c_p * deltaT / L_f)
//
// With the steel values in this package (c_p = 0.80 kJ/(kg K),
// L_f = 270 kJ/kg) a 50 K superheat extends the freezing time by roughly
// 15%. A zero superheat reproduces the pure Chvorinov result exactly.
package chvorinov

// SuperheatFactor returns the multiplier applied to the pure Chvorinov
// freezing time when the metal is superheated by deltaT kelvin above the
// liquidus. The heat capacity and the latent heat are supplied so that
// materials other than steel can be handled without new constants.
func SuperheatFactor(deltaT, heatCapacity, latentHeat float64) float64 {
	if deltaT <= 0 {
		return 1.0
	}
	if latentHeat <= 0 {
		return 1.0
	}
	return 1.0 + heatCapacity*deltaT/latentHeat
}

// ApplySuperheat multiplies a pure freezing time by the superheat factor.
// It is a pure function: the input time is never modified in place.
func ApplySuperheat(tf, deltaT, heatCapacity, latentHeat float64) float64 {
	return tf * SuperheatFactor(deltaT, heatCapacity, latentHeat)
}

// SuperheatExtension returns the absolute amount of time added by the
// superheat correction, tf' - tf. Negative values are clamped to zero so
// the report can print "0.000 min added" for pure Chvorinov cases.
func SuperheatExtension(tf, deltaT, heatCapacity, latentHeat float64) float64 {
	ext := tf * (SuperheatFactor(deltaT, heatCapacity, latentHeat) - 1.0)
	if ext < 0 {
		return 0.0
	}
	return ext
}

// SteelSuperheatFactor is the superheat factor using the steel properties
// stored in this package. It exists so the CLI does not have to pass the
// material constants by hand.
func SteelSuperheatFactor(deltaT float64) float64 {
	return SuperheatFactor(deltaT, SteelHeatCapacity, SteelLatentHeat)
}

// SteelSuperheatExtension is the absolute extension using steel constants.
func SteelSuperheatExtension(tf, deltaT float64) float64 {
	return SuperheatExtension(tf, deltaT, SteelHeatCapacity, SteelLatentHeat)
}
