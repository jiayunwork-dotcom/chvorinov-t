package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"chvorinov-t/internal/chvorinov"
	"chvorinov-t/internal/cli"
	"chvorinov-t/internal/geometry"
)

type shapeRowView struct {
	Label        string  `json:"label"`
	SizeDesc     string  `json:"size_desc"`
	Volume       float64 `json:"volume"`
	Area         float64 `json:"area"`
	Modulus      float64 `json:"modulus"`
	Time         float64 `json:"time"`
	IdealModulus float64 `json:"ideal_modulus"`
	IdealFormula string  `json:"ideal_formula"`
}

type riserView struct {
	CastingModulus  float64 `json:"casting_modulus"`
	RiserModulus    float64 `json:"riser_modulus"`
	RequiredModulus float64 `json:"required_modulus"`
	Status          string  `json:"status"`
	Warning         string  `json:"warning"`
}

type endpointDoc struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Doc    string `json:"doc"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]string{
		"service": "chvorinov-t",
		"freeze":  "POST /api/freeze",
		"compare": "POST /api/compare",
		"scale":   "POST /api/scale",
		"riser":   "POST /api/riser",
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func handleFreeze(w http.ResponseWriter, r *http.Request) {
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
	table := chvorinov.CompareShapes(shape.V, in.C, n, 0)
	var riser *riserView
	if in.Riser != nil {
		rm, ok, rerr := in.RiserModulus()
		if rerr != nil {
			writeError(w, rerr)
			return
		}
		if ok {
			rc := chvorinov.CheckRiser(res.Modulus, rm)
			riser = &riserView{
				CastingModulus:  rc.CastingModulus,
				RiserModulus:    rc.RiserModulus,
				RequiredModulus: rc.RequiredModulus,
				Status:          rc.Status.String(),
				Warning:         rc.Warning,
			}
		}
	}
	writeJSON(w, map[string]any{
		"label":             in.Label,
		"shape":             shape.Kind.String(),
		"volume":            res.Volume,
		"area":              res.Area,
		"modulus":           res.Modulus,
		"mold_const":        res.MoldConst,
		"exponent":          res.Exponent,
		"freeze_time_min":   res.FreezeTime,
		"superheat_k":       res.SuperheatK,
		"superheat_applied": res.SuperheatApplied,
		"min_area":          res.MinArea,
		"table":             shapeRows(table),
		"riser":             riser,
	})
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
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
	table := chvorinov.CompareShapes(shape.V, in.C, n, 0)
	writeJSON(w, map[string]any{
		"volume":     table.Volume,
		"mold_const": table.MoldConst,
		"exponent":   table.Exponent,
		"slowest":    table.Slowest,
		"fastest":    table.Fastest,
		"rows":       shapeRows(table),
	})
}

func handleScale(w http.ResponseWriter, r *http.Request) {
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
	rep := chvorinov.CrossScale(res.Modulus, res.FreezeTime, 2, n)
	writeJSON(w, map[string]any{
		"factor":         rep.Factor,
		"exponent":       rep.Exponent,
		"ref_modulus":    rep.RefModulus,
		"scaled_modulus": rep.ScaledModulus,
		"ref_time":       rep.RefTime,
		"scaled_time":    rep.ScaledTime,
		"time_ratio":     rep.TimeRatio,
	})
}

func handleRiser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	in, shape, err := decodeCasting(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if in.Riser == nil {
		writeError(w, errNoRiser)
		return
	}
	rm, ok, err := in.RiserModulus()
	if err != nil {
		writeError(w, err)
		return
	}
	if !ok {
		writeError(w, errNoRiser)
		return
	}
	n := in.Exponent()
	res, err := chvorinov.Compute(shape.V, shape.A, in.C, n, in.SuperheatK)
	if err != nil {
		writeError(w, err)
		return
	}
	rc := chvorinov.CheckRiser(res.Modulus, rm)
	writeJSON(w, riserView{
		CastingModulus:  rc.CastingModulus,
		RiserModulus:    rc.RiserModulus,
		RequiredModulus: rc.RequiredModulus,
		Status:          rc.Status.String(),
		Warning:         rc.Warning,
	})
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, _, err := decodeCasting(r); err != nil {
		writeJSON(w, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]any{
		"name":    "chvorinov-t",
		"version": "1.0.0",
		"routes":  []string{"/api/freeze", "/api/compare", "/api/scale", "/api/riser", "/api/validate", "/api/rules", "/api/modulus", "/api/superheat", "/health"},
	})
}

func handleEndpoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, []endpointDoc{
		{Method: "POST", Path: "/api/freeze", Doc: "凝固时间核算"},
		{Method: "POST", Path: "/api/compare", Doc: "同体积形状对照"},
		{Method: "POST", Path: "/api/scale", Doc: "相似缩放交叉规则"},
		{Method: "POST", Path: "/api/riser", Doc: "冒口模数检查"},
		{Method: "POST", Path: "/api/validate", Doc: "输入校验"},
		{Method: "POST", Path: "/api/rules", Doc: "交叉规则核验"},
		{Method: "POST", Path: "/api/modulus", Doc: "模数与理想模数"},
		{Method: "POST", Path: "/api/superheat", Doc: "过热时间修正"},
		{Method: "GET", Path: "/health", Doc: "健康检查"},
	})
}

func decodeCasting(r *http.Request) (*cli.Input, geometry.Shape, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, geometry.Shape{}, err
	}
	in, err := cli.DecodeInput(bytes.NewReader(body))
	if err != nil {
		return nil, geometry.Shape{}, err
	}
	shape, err := in.ToGeometry()
	if err != nil {
		return nil, geometry.Shape{}, err
	}
	return in, shape, nil
}

func shapeRows(t chvorinov.ShapeTable) []shapeRowView {
	rows := make([]shapeRowView, 0, len(t.Rows))
	for _, r := range t.Rows {
		rows = append(rows, shapeRowView{
			Label:        r.Label,
			SizeDesc:     r.SizeDesc,
			Volume:       r.Volume,
			Area:         r.Area,
			Modulus:      r.Modulus,
			Time:         r.Time,
			IdealModulus: r.IdealModulus,
			IdealFormula: r.IdealFormula,
		})
	}
	return rows
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
