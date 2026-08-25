package chvorinov

import "fmt"

type RiserStatus int

const (
	RiserOk RiserStatus = iota
	RiserTooSmall
	RiserEqual
)

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

type RiserCheck struct {
	CastingModulus  float64
	RiserModulus    float64
	RequiredModulus float64
	Status          RiserStatus
	Warning         string
}

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

func (c RiserCheck) Describe() string {
	return fmt.Sprintf(
		"riser M=%.4f cm  casting M=%.4f cm  required=%.4f cm  status=%s",
		c.RiserModulus, c.CastingModulus, c.RequiredModulus, c.Status,
	)
}

func (c RiserCheck) HasWarning() bool {
	return c.Warning != ""
}

func MinRiserModulus(castingModulus float64) float64 {
	return castingModulus * RiserSafetyFactor
}
