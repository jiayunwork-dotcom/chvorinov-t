package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"chvorinov-t/internal/chvorinov"
	"chvorinov-t/internal/geometry"
)

type RiserInput struct {
	Shape     string   `json:"shape,omitempty"`
	Edge      *float64 `json:"edge,omitempty"`
	Thickness *float64 `json:"thickness,omitempty"`
	Width     *float64 `json:"width,omitempty"`
	Length    *float64 `json:"length,omitempty"`
	Radius    *float64 `json:"radius,omitempty"`
	Height    *float64 `json:"height,omitempty"`
	Volume    *float64 `json:"volume,omitempty"`
	Area      *float64 `json:"area,omitempty"`
	C         *float64 `json:"C,omitempty"`
	N         *float64 `json:"n,omitempty"`
}

type Input struct {
	Label      string      `json:"label,omitempty"`
	Shape      string      `json:"shape,omitempty"`
	Edge       *float64    `json:"edge,omitempty"`
	Thickness  *float64    `json:"thickness,omitempty"`
	Width      *float64    `json:"width,omitempty"`
	Length     *float64    `json:"length,omitempty"`
	Radius     *float64    `json:"radius,omitempty"`
	Height     *float64    `json:"height,omitempty"`
	Volume     *float64    `json:"volume,omitempty"`
	Area       *float64    `json:"area,omitempty"`
	C          float64     `json:"C"`
	N          *float64    `json:"n,omitempty"`
	SuperheatK float64     `json:"superheat_k,omitempty"`
	Riser      *RiserInput `json:"riser,omitempty"`
}

func dimsFrom(shape string, e, t, w, l, r, h *float64) geometry.Dims {
	return geometry.Dims{
		Edge:      valueOf(e),
		Thickness: valueOf(t),
		Width:     valueOf(w),
		Length:    valueOf(l),
		Radius:    valueOf(r),
		Height:    valueOf(h),
	}
}

func valueOf(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func (in *Input) exponent() float64 {
	if in.N != nil {
		return *in.N
	}
	return chvorinov.DefaultExponent
}

func (in *Input) geometryShape() (geometry.Shape, error) {
	kind, ok := geometry.ParseShapeKind(in.Shape)
	if !ok {
		return geometry.Shape{}, fmt.Errorf("unknown shape %q (want cube, plate, cylinder, sphere, or custom)", in.Shape)
	}
	if kind == geometry.ShapeCustom {
		if in.Volume == nil || in.Area == nil {
			return geometry.Shape{}, errors.New("custom shape requires both \"volume\" and \"area\" fields")
		}
		return geometry.FromVolumeArea(*in.Volume, *in.Area), nil
	}

	d := dimsFrom(in.Shape, in.Edge, in.Thickness, in.Width, in.Length, in.Radius, in.Height)
	if err := checkDims(kind, d); err != nil {
		return geometry.Shape{}, err
	}
	return geometry.New(kind, d), nil
}

func checkDims(kind geometry.ShapeKind, d geometry.Dims) error {
	switch kind {
	case geometry.ShapeCube:
		if d.Edge <= 0 {
			return errors.New("cube requires a positive \"edge\"")
		}
	case geometry.ShapePlate:
		if d.Thickness <= 0 || d.Width <= 0 || d.Length <= 0 {
			return errors.New("plate requires positive \"thickness\", \"width\", and \"length\"")
		}
	case geometry.ShapeCylinder:
		if d.Radius <= 0 || d.Height <= 0 {
			return errors.New("cylinder requires positive \"radius\" and \"height\"")
		}
	case geometry.ShapeSphere:
		if d.Radius <= 0 {
			return errors.New("sphere requires a positive \"radius\"")
		}
	}
	return nil
}

func (in *Input) riserModulus() (float64, bool, error) {
	if in.Riser == nil {
		return 0, false, nil
	}
	r := in.Riser
	kind, ok := geometry.ParseShapeKind(r.Shape)
	if !ok {
		return 0, true, fmt.Errorf("unknown riser shape %q", r.Shape)
	}
	var m float64
	if kind == geometry.ShapeCustom {
		if r.Volume == nil || r.Area == nil {
			return 0, true, errors.New("custom riser requires both \"volume\" and \"area\" fields")
		}
		m = geometry.Modulus(*r.Volume, *r.Area)
	} else {
		d := dimsFrom(r.Shape, r.Edge, r.Thickness, r.Width, r.Length, r.Radius, r.Height)
		if err := checkDims(kind, d); err != nil {
			return 0, true, err
		}
		m = geometry.New(kind, d).Modulus()
	}
	return m, true, nil
}

func LoadInput(path string) (*Input, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	return DecodeInput(f)
}

func DecodeInput(r io.Reader) (*Input, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var in Input
	if err := dec.Decode(&in); err != nil {
		return nil, fmt.Errorf("decode input: %w", err)
	}
	if in.C <= 0 {
		return nil, fmt.Errorf("mold constant C must be strictly positive (got %v)", in.C)
	}
	if in.SuperheatK < 0 {
		return nil, fmt.Errorf("superheat_k must be >= 0 (got %v)", in.SuperheatK)
	}
	if _, err := in.geometryShape(); err != nil {
		return nil, err
	}
	return &in, nil
}

func (in *Input) ToGeometry() (geometry.Shape, error) {
	return in.geometryShape()
}

func (in *Input) RiserModulus() (float64, bool, error) {
	return in.riserModulus()
}

func (in *Input) Exponent() float64 {
	return in.exponent()
}
