// This file implements the riser check. A riser is a reservoir of liquid
// metal attached to the casting; it must stay liquid until the casting has
// finished freezing, otherwise it cannot feed the shrinkage as the casting
// solidifies. Under Chvorinov's rule the riser must therefore have a
// larger modulus than the casting (the riser "freezes last"). When the
// riser modulus does not exceed the casting modulus, the tool reports a
// warning rather than failing, because the casting itself is still
// geometrically valid.
package chvorinov

import "fmt"

// RiserStatus describes the outcome of a riser comparison.
type RiserStatus int

const (
	// RiserOk means the riser modulus is large enough and the riser
	// solidifies after the casting.
	RiserOk RiserStatus = iota
	// RiserTooSmall means the riser modulus is not larger than the casting
	// modulus; the riser would freeze first and fail to feed.
	RiserTooSmall
	// RiserEqual means the riser and casting moduli are equal within
	// tolerance, which gives no guarantee of feeding order.
	RiserEqual
)

// String renders the status for the report.
func (s RiserStatus) String() string {
	switch s {
	case RiserOk:
		return "ok"
	case RiserEqual:
		return "equal"
	default:
		return "too-small"
	}
}

// RiserCheck is the result of comparing a riser modulus against a casting
// modulus. The required modulus is the casting modulus times the safety
// factor; the margin is how far above that requirement the riser sits.
type RiserCheck struct {
	CastingModulus  float64
	RiserModulus    float64
	RequiredModulus float64
	Status          RiserStatus
	Warning         string
}

// CheckRiser compares a riser modulus against a casting modulus. A riser
// is adequate only when its modulus exceeds the required modulus (the
// casting modulus scaled by RiserSafetyFactor). The check never fails the
// whole computation; it returns a warning when the riser is inadequate.
func CheckRiser(castingModulus, riserModulus float64) RiserCheck {
	required := castingModulus * RiserSafetyFactor
	status := RiserTooSmall
	msg := "riser modulus must exceed the casting modulus to solidify last"
	if riserModulus > required {
		status = RiserOk
		msg = ""
	} else if NearlyEqualRatio(riserModulus, required) {
		status = RiserEqual
		msg = "riser modulus equals the casting modulus; no feeding margin"
	}
	return RiserCheck{
		CastingModulus:  castingModulus,
		RiserModulus:    riserModulus,
		RequiredModulus: required,
		Status:          status,
		Warning:         msg,
	}
}

// NearlyEqualRatio reports whether two positive values are equal within
// ZeroTolerance relative to the larger one.
func NearlyEqualRatio(x, y float64) bool {
	larger := x
	if y > larger {
		larger = y
	}
	if larger <= 0 {
		return x == y
	}
	diff := x - y
	if diff < 0 {
		diff = -diff
	}
	return diff/larger <= ZeroTolerance
}

// Describe renders the riser check as a single human-readable line.
func (c RiserCheck) Describe() string {
	return fmt.Sprintf(
		"riser M=%.4f cm  casting M=%.4f cm  required=%.4f cm  status=%s",
		c.RiserModulus, c.CastingModulus, c.RequiredModulus, c.Status,
	)
}

// HasWarning reports whether the check produced a warning for the user.
func (c RiserCheck) HasWarning() bool {
	return c.Warning != ""
}

// MinRiserModulus returns the smallest riser modulus that passes the
// check for a given casting modulus.
func MinRiserModulus(castingModulus float64) float64 {
	return castingModulus * RiserSafetyFactor
}
