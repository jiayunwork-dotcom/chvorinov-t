// Package chvorinov implements the casting solidification kernel based on
// Chvorinov's rule:
//
//	tf = C * (V/A)^n = C * M^n
//
// where V is the casting volume, A its surface area, M = V/A the
// solidification modulus, C the mold constant, and n the exponent
// (commonly close to 2). The package validates the inputs, computes the
// freezing time, checks the riser, applies an optional superheat
// correction, and verifies the cross rules (similarity scaling, shape
// comparison, physical feasibility) that the rest of the tool relies on.
package chvorinov

// DefaultExponent is the exponent used when the input file does not set n
// explicitly. Casting practice places n close to 2 for most molds.
const DefaultExponent = 2.0

// RiserSafetyFactor is the multiplier applied to the casting modulus to
// obtain the minimum riser modulus that guarantees the riser solidifies
// after the casting. A value above 1 gives the riser a margin over the
// "must be larger" condition.
const RiserSafetyFactor = 1.0

// DefaultRiserMargin is the fraction the riser modulus must exceed the
// casting modulus by to pass the riser check without a warning. The
// default keeps the strict Chvorinov condition: strictly greater is
// enough.
const DefaultRiserMargin = 0.0

// ZeroTolerance is the relative tolerance used when deciding whether a
// ratio is exactly one (e.g. whether the riser modulus equals the casting
// modulus). Inputs are user-provided decimals, so exact equality cannot
// be assumed.
const ZeroTolerance = 1e-12

// SteelHeatCapacity is the specific heat capacity of liquid steel in
// kJ/(kg K), used by the optional superheat correction.
const SteelHeatCapacity = 0.80

// SteelLatentHeat is the latent heat of fusion of steel in kJ/kg.
const SteelLatentHeat = 270.0

// SteelDensity is the density of steel in kg/cm^3, used only for display
// of mass-related estimates in reports.
const SteelDensity = 0.00785

// DefaultMoldConstantSteel is a typical Chvorinov constant for green-sand
// molds and steel castings in min/cm^2 when n = 2. It is provided so the
// example files and the README can quote a representative value; it is not
// applied automatically.
const DefaultMoldConstantSteel = 3.0

// MinExponentFloor is the smallest exponent accepted by the validator.
// n <= 0 is rejected because a non-positive exponent would make the
// freezing time shrink with the modulus.
const MinExponentFloor = 0.0

// MaxPrintPrecision is the number of decimal places used when rendering
// numbers in reports.
const MaxPrintPrecision = 3
