package server

import (
	"net/http"

	"chvorinov-t/internal/chvorinov"
)

type ruleView struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"`
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	in, shape, err := decodeCasting(r)
	if err != nil {
		writeError(w, err)
		return
	}
	rules := chvorinov.VerifyCrossRules(shape.V, in.C, in.Exponent())
	rows := make([]ruleView, 0, len(rules.Rules))
	for _, rr := range rules.Rules {
		rows = append(rows, ruleView{Name: rr.Name, Pass: rr.Pass, Detail: rr.Detail})
	}
	writeJSON(w, map[string]any{
		"volume":     rules.Volume,
		"mold_const": rules.MoldConst,
		"exponent":   rules.Exponent,
		"all_pass":   rules.AllPass,
		"rules":      rows,
	})
}

func handleModulus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	in, shape, err := decodeCasting(r)
	if err != nil {
		writeError(w, err)
		return
	}
	rep := chvorinov.NewModulusReport(in.Label, shape.V, shape.A, shape.Kind, shape.Dims)
	writeJSON(w, map[string]any{
		"shape_label":       rep.ShapeLabel,
		"volume":            rep.Volume,
		"area":              rep.Area,
		"modulus":           rep.Modulus,
		"min_area":          rep.MinArea,
		"ideal_modulus":     rep.IdealModulus,
		"surface_to_volume": chvorinov.SurfaceToVolume(rep.Volume, rep.Area),
	})
}

func handleSuperheat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	in, shape, err := decodeCasting(r)
	if err != nil {
		writeError(w, err)
		return
	}
	n := in.Exponent()
	res, err := chvorinov.Compute(shape.V, shape.A, in.C, n, in.SuperheatK)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{
		"superheat_k":          in.SuperheatK,
		"factor":               chvorinov.SteelSuperheatFactor(in.SuperheatK),
		"extension_min":        chvorinov.SteelSuperheatExtension(res.FreezeTime, in.SuperheatK),
		"base_freeze_time_min": res.FreezeTime,
	})
}
