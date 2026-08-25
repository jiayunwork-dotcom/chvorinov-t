package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"chvorinov-t/internal/geometry"
)

func writeTempJSON(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestLoadInputCube(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3,"n":2}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	if in.C != 3 || in.exponent() != 2 {
		t.Errorf("C=%v n=%v, want C=3 n=2", in.C, in.exponent())
	}
	s, err := in.geometryShape()
	if err != nil {
		t.Fatalf("geometryShape = %v", err)
	}
	if s.Kind != geometry.ShapeCube {
		t.Errorf("kind = %v, want cube", s.Kind)
	}
	if got, want := s.V, 1000.0; got != want {
		t.Errorf("volume = %v, want %v", got, want)
	}
	if got, want := s.A, 600.0; got != want {
		t.Errorf("area = %v, want %v", got, want)
	}
}

func TestLoadInputDefaultExponent(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	if in.exponent() != 2.0 {
		t.Errorf("exponent() = %v, want 2.0", in.exponent())
	}
}

func TestLoadInputCustomVolumeArea(t *testing.T) {
	p := writeTempJSON(t, "custom.json", `{"shape":"custom","volume":1000,"area":600,"C":3}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	s, err := in.geometryShape()
	if err != nil {
		t.Fatalf("geometryShape = %v", err)
	}
	if s.Kind != geometry.ShapeCustom {
		t.Errorf("kind = %v, want custom", s.Kind)
	}
	if got, want := s.Modulus(), 1000.0/600.0; got != want {
		t.Errorf("modulus = %v, want %v", got, want)
	}
}

func TestLoadInputCustomRequiresVolumeArea(t *testing.T) {
	p := writeTempJSON(t, "custom.json", `{"shape":"custom","volume":1000,"C":3}`)
	if _, err := LoadInput(p); err == nil {
		t.Fatal("LoadInput(custom without area) = nil, want error")
	}
}

func TestLoadInputUnknownShape(t *testing.T) {
	p := writeTempJSON(t, "bad.json", `{"shape":"torus","C":3,"n":2}`)
	_, err := LoadInput(p)
	if err == nil {
		t.Fatal("LoadInput(unknown shape) = nil, want error")
	}
	if !strings.Contains(err.Error(), "unknown shape") {
		t.Errorf("error %q should mention the unknown shape", err)
	}
}

func TestLoadInputMissingMoldConst(t *testing.T) {
	p := writeTempJSON(t, "bad.json", `{"shape":"cube","edge":10,"n":2}`)
	if _, err := LoadInput(p); err == nil {
		t.Fatal("LoadInput(no C) = nil, want error")
	}
}

func TestLoadInputRejectsUnknownField(t *testing.T) {
	p := writeTempJSON(t, "bad.json", `{"shape":"cube","edge":10,"C":3,"voluem":1000}`)
	if _, err := LoadInput(p); err == nil {
		t.Fatal("LoadInput(unknown field) = nil, want error")
	}
}

func TestLoadInputCubeMissingEdge(t *testing.T) {
	p := writeTempJSON(t, "bad.json", `{"shape":"cube","C":3}`)
	if _, err := LoadInput(p); err == nil {
		t.Fatal("LoadInput(cube without edge) = nil, want error")
	}
}

func TestLoadInputSphereRadius(t *testing.T) {
	p := writeTempJSON(t, "sphere.json", `{"shape":"sphere","radius":6.204,"C":3,"n":2}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	s, err := in.geometryShape()
	if err != nil {
		t.Fatalf("geometryShape = %v", err)
	}
	if s.Kind != geometry.ShapeSphere {
		t.Errorf("kind = %v, want sphere", s.Kind)
	}
}

func TestLoadInputPlate(t *testing.T) {
	p := writeTempJSON(t, "plate.json", `{"shape":"plate","thickness":2.5,"width":20,"length":20,"C":3}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	s, err := in.geometryShape()
	if err != nil {
		t.Fatalf("geometryShape = %v", err)
	}
	if s.Kind != geometry.ShapePlate {
		t.Errorf("kind = %v, want plate", s.Kind)
	}
}

func TestRiserModulusComputesCylinder(t *testing.T) {
	p := writeTempJSON(t, "riser.json", `{
		"shape":"cube","edge":10,"C":3,
		"riser":{"shape":"cylinder","radius":6,"height":30}
	}`)
	in, err := LoadInput(p)
	if err != nil {
		t.Fatalf("LoadInput = %v", err)
	}
	m, ok, err := in.riserModulus()
	if err != nil {
		t.Fatalf("riserModulus = %v", err)
	}
	if !ok {
		t.Fatal("riserModulus reported no riser")
	}
	if m <= 0 {
		t.Errorf("riser modulus = %v, want > 0", m)
	}
}
