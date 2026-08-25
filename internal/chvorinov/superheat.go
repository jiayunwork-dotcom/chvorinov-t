package chvorinov

func SuperheatFactor(deltaT, heatCapacity, latentHeat float64) float64 {
	if deltaT <= 0 {
		return 1.0
	}
	if latentHeat <= 0 {
		return 1.0
	}
	return 1.0 + heatCapacity*deltaT/latentHeat
}

func ApplySuperheat(tf, deltaT, heatCapacity, latentHeat float64) float64 {
	return tf * SuperheatFactor(deltaT, heatCapacity, latentHeat)
}

func SuperheatExtension(tf, deltaT, heatCapacity, latentHeat float64) float64 {
	ext := tf * (SuperheatFactor(deltaT, heatCapacity, latentHeat) - 1.0)
	if ext < 0 {
		return 0.0
	}
	return ext
}

func SteelSuperheatFactor(deltaT float64) float64 {
	return SuperheatFactor(deltaT, SteelHeatCapacity, SteelLatentHeat)
}

func SteelSuperheatExtension(tf, deltaT float64) float64 {
	return SuperheatExtension(tf, deltaT, SteelHeatCapacity, SteelLatentHeat)
}
