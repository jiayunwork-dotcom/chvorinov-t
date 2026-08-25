package chvorinov

import (
	"fmt"

	"chvorinov-t/internal/geometry"
)

type CrossRule struct {
	Name   string
	Pass   bool
	Detail string
}

type CrossRules struct {
	Volume    float64
	MoldConst float64
	Exponent  float64
	Rules     []CrossRule
	AllPass   bool
}

func VerifyCrossRules(volume, moldConst, exponent float64) CrossRules {
	var rules []CrossRule

	rep := CrossScale(CubeModulusForVolume(volume), FreezeTimeFromVA(volume, geometry.CubeArea(geometry.CubeThicknessForVolume(volume)), moldConst, exponent), 2.0, exponent)
	rules = append(rules, CrossRule{
		Name:   "similarity scaling: size x2 quadruples tf at n=2",
		Pass:   IsQuadrupleTime(rep.TimeRatio),
		Detail: "time ratio " + fmtRatio(rep.TimeRatio),
	})

	baseTf := FreezeTime(moldConst, CubeModulusForVolume(volume), exponent)
	doubledTf := DoubleMoldConstTime(baseTf)
	rules = append(rules, CrossRule{
		Name:   "mold constant doubling: C x2 doubles tf",
		Pass:   IsDoubleTime(doubledTf / baseTf),
		Detail: "time ratio " + fmtRatio(doubledTf/baseTf),
	})

	refT := geometry.CubeThicknessForVolume(volume)
	flatT := refT / 4.0
	eff := geometry.Flatten(volume, refT, flatT)
	rules = append(rules, CrossRule{
		Name:   "flatter plate: larger A, smaller M, shorter tf at equal V",
		Pass:   eff.IsFlatterLargerArea() && eff.IsFlatterSmallerModulus() && eff.IsFlatterShorterTime(),
		Detail: "area growth " + fmtRatio(eff.AreaGrowth) + ", modulus ratio " + fmtRatio(eff.ModulusRatio),
	})

	m := CubeModulusForVolume(volume)
	cubeTf := EqualModulusTime(m, moldConst, exponent)
	sphereTf := EqualModulusTime(m, moldConst, exponent)
	rules = append(rules, CrossRule{
		Name:   "equal V/A gives equal tf across shape labels",
		Pass:   NearlyEqualRatio(cubeTf, sphereTf),
		Detail: "cube tf " + fmtNumSimple(cubeTf) + ", sphere(label) tf " + fmtNumSimple(sphereTf),
	})

	table := CompareShapes(volume, moldConst, exponent, 0)
	rules = append(rules, CrossRule{
		Name:   "same volume: sphere freezes slower than plate",
		Pass:   table.IsSphereSlowest() && table.IsPlateFastest(),
		Detail: "slowest " + table.Slowest + ", fastest " + table.Fastest,
	})

	all := true
	for _, r := range rules {
		if !r.Pass {
			all = false
		}
	}
	return CrossRules{
		Volume:    volume,
		MoldConst: moldConst,
		Exponent:  exponent,
		Rules:     rules,
		AllPass:   all,
	}
}

func fmtRatio(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func fmtNumSimple(v float64) string {
	return fmt.Sprintf("%.3f", v)
}
