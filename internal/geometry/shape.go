// Package geometry provides the shape geometry used by casting
// solidification calculations: volume, surface area, and the
// solidification modulus M = V/A for the standard shape classes
// (cube, plate, cylinder, sphere) plus a custom class that takes
// raw volume and area directly from the user.
package geometry

// ShapeKind identifies a casting shape class. The kind selects which
// dimensional formula is used to compute volume and surface area.
type ShapeKind int

const (
	// ShapeCustom uses explicitly supplied volume and area; no dimension
	// formula is applied.
	ShapeCustom ShapeKind = iota
	// ShapeCube is a rectangular block with all three edges equal.
	ShapeCube
	// ShapePlate is a flat slab whose thickness is small compared with
	// its width and length.
	ShapePlate
	// ShapeCylinder is a right circular cylinder.
	ShapeCylinder
	// ShapeSphere is a solid ball.
	ShapeSphere
)

// String returns a stable, human-readable label for the shape class.
func (k ShapeKind) String() string {
	switch k {
	case ShapeCube:
		return "cube"
	case ShapePlate:
		return "plate"
	case ShapeCylinder:
		return "cylinder"
	case ShapeSphere:
		return "sphere"
	default:
		return "custom"
	}
}

// ParseShapeKind maps a shape label to its ShapeKind. An empty label or
// the label "custom" maps to ShapeCustom so callers can fall back to raw
// volume and area. The second return value reports whether the label was
// recognised at all.
func ParseShapeKind(label string) (ShapeKind, bool) {
	switch label {
	case "cube":
		return ShapeCube, true
	case "plate":
		return ShapePlate, true
	case "cylinder":
		return ShapeCylinder, true
	case "sphere":
		return ShapeSphere, true
	case "", "custom":
		return ShapeCustom, true
	default:
		return ShapeCustom, false
	}
}

// Dims carries the linear dimensions of a shape. Only the fields relevant
// to a particular ShapeKind are interpreted; the others are ignored.
type Dims struct {
	Edge      float64 // cube side length
	Thickness float64 // plate thickness
	Width     float64 // plate width
	Length    float64 // plate length
	Radius    float64 // cylinder or sphere radius
	Height    float64 // cylinder height
}

// Shape is a concrete geometry object with a known volume and surface
// area. The kind and dimensions are kept alongside V and A so that the
// modulus and the freeze time can be derived later without re-reading the
// input file.
type Shape struct {
	Kind ShapeKind
	Dims Dims
	V    float64 // volume in cm^3
	A    float64 // surface area in cm^2
}

// New builds a Shape from a kind and its dimensions. It computes the
// exact volume and surface area for the shape class. An empty Dims for a
// recognised kind yields a zero-volume shape, which callers must reject
// during validation.
func New(kind ShapeKind, d Dims) Shape {
	switch kind {
	case ShapeCube:
		return Shape{
			Kind: ShapeCube,
			Dims: d,
			V:    CubeVolume(d.Edge),
			A:    CubeArea(d.Edge),
		}
	case ShapePlate:
		return Shape{
			Kind: ShapePlate,
			Dims: d,
			V:    PlateVolume(d.Thickness, d.Width, d.Length),
			A:    PlateArea(d.Thickness, d.Width, d.Length),
		}
	case ShapeCylinder:
		return Shape{
			Kind: ShapeCylinder,
			Dims: d,
			V:    CylinderVolume(d.Radius, d.Height),
			A:    CylinderArea(d.Radius, d.Height),
		}
	case ShapeSphere:
		return Shape{
			Kind: ShapeSphere,
			Dims: d,
			V:    SphereVolume(d.Radius),
			A:    SphereArea(d.Radius),
		}
	default:
		return Shape{Kind: ShapeCustom, Dims: d}
	}
}

// FromVolumeArea builds a Shape from an explicit volume and area. The
// shape kind is recorded as custom because no dimensional formula applies.
func FromVolumeArea(v, a float64) Shape {
	return Shape{Kind: ShapeCustom, V: v, A: a}
}

// Modulus returns the solidification modulus M = V/A. Validation of the
// inputs is the caller's responsibility; the function itself does not
// guard against a zero area because division happens only after the
// physical checks have passed.
func (s Shape) Modulus() float64 {
	return s.V / s.A
}
