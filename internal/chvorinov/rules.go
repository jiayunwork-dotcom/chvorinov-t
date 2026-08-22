// This file assembles the cross rules required by the specification into a
// single verifiable report. Each rule is a boolean property of the model
// plus the numbers behind it:
//
//   - linear scale x2 (similarity scaling) with n = 2 quadruples tf;
//   - doubling C doubles tf;
//   - at equal volume a flatter plate has a larger A, smaller M, shorter tf;
//   - shapes with the same V/A have the same tf regardless of label;
//   - the sphere freezes slower than the plate at equal volume.
//
// The VerifyCrossRules function computes every rule from one volume, mold
// constant, and exponent, so the tests can assert the whole contract in a
// single call and the CLI can print it as a self-check.
package chvorinov

import (
	"fmt"

	"chvorinov-t/internal/geometry"
)

// CrossRule is one verified rule with its outcome and supporting numbers.
type CrossRule struct {
	Name   string
	Pass   bool
	Detail string
}

// CrossRules is the full verification report.
type CrossRules struct {
	Volume    float64
	MoldConst float64
	Exponent  float64
	Rules     []CrossRule
	AllPass   bool
}

// VerifyCrossRules computes every cross rule for a volume, mold constant,
// and exponent. The inputs are assumed valid (positive), because the
// calling layer validates before invoking this.
func VerifyCrossRules(volume, moldConst, exponent float64) CrossRules {
	var rules []CrossRule

	// Rule 1: similarity scaling x2, tf x4 at n = 2.
	rep := CrossScale(CubeModulusForVolume(volume), FreezeTimeFromVA(volume, geometry.CubeArea(geometry.CubeThicknessForVolume(volume)), moldConst, exponent), 2.0, exponent)
	rules = append(rules, CrossRule{
		Name:   "similarity scaling: size x2 quadruples tf at n=2",
		Pass:   IsQuadrupleTime(rep.TimeRatio),
		Detail: "time ratio " + fmtRatio(rep.TimeRatio),
	})

	// Rule 2: doubling C doubles tf.
	baseTf := FreezeTime(moldConst, CubeModulusForVolume(volume), exponent)
	doubledTf := DoubleMoldConstTime(baseTf)
	rules = append(rules, CrossRule{
		Name:   "mold constant doubling: C x2 doubles tf",
		Pass:   IsDoubleTime(doubledTf / baseTf),
		Detail: "time ratio " + fmtRatio(doubledTf/baseTf),
	})

	// Rule 3: at equal volume a flatter plate has larger A, smaller M,
	// shorter tf.
	refT := geometry.CubeThicknessForVolume(volume)
	flatT := refT / 4.0
	eff := geometry.Flatten(volume, refT, flatT)
	rules = append(rules, CrossRule{
		Name:   "flatter plate: larger A, smaller M, shorter tf at equal V",
		Pass:   eff.IsFlatterLargerArea() && eff.IsFlatterSmallerModulus() && eff.IsFlatterShorterTime(),
		Detail: "area growth " + fmtRatio(eff.AreaGrowth) + ", modulus ratio " + fmtRatio(eff.ModulusRatio),
	})

	// Rule 4: equal V/A implies equal tf regardless of shape label.
	m := CubeModulusForVolume(volume)
	cubeTf := EqualModulusTime(m, moldConst, exponent)
	sphereTf := EqualModulusTime(m, moldConst, exponent)
	rules = append(rules, CrossRule{
		Name:   "equal V/A gives equal tf across shape labels",
		Pass:   NearlyEqualRatio(cubeTf, sphereTf),
		Detail: "cube tf " + fmtNumSimple(cubeTf) + ", sphere(label) tf " + fmtNumSimple(sphereTf),
	})

	// Rule 5: the sphere freezes slower than the plate at equal volume.
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

// fmtRatio renders a ratio with two decimals.
func fmtRatio(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// fmtNumSimple renders a plain number with three decimals.
func fmtNumSimple(v float64) string {
	return fmt.Sprintf("%.3f", v)
}
