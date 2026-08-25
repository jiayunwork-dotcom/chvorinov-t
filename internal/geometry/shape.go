package geometry

type ShapeKind int

const (
	ShapeCustom ShapeKind = iota
	ShapeCube
	ShapePlate
	ShapeCylinder
	ShapeSphere
)

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

type Dims struct {
	Edge      float64
	Thickness float64
	Width     float64
	Length    float64
	Radius    float64
	Height    float64
}

type Shape struct {
	Kind ShapeKind
	Dims Dims
	V    float64
	A    float64
}

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

func FromVolumeArea(v, a float64) Shape {
	return Shape{Kind: ShapeCustom, V: v, A: a}
}

func (s Shape) Modulus() float64 {
	return s.V / s.A
}
