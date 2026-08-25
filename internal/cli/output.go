package cli

import (
	"fmt"
	"io"
	"strings"

	"chvorinov-t/internal/chvorinov"
)

func fmtNum(v float64) string {
	return fmt.Sprintf("%.3f", v)
}

func fmtMod(v float64) string {
	return fmt.Sprintf("%.4f", v)
}

func writeHeading(w io.Writer, title string) {
	fmt.Fprintf(w, "\n%s\n%s\n", title, strings.Repeat("-", len(title)))
}

func writeResult(w io.Writer, r chvorinov.Result, label string) {
	writeHeading(w, "casting")
	if label != "" {
		fmt.Fprintf(w, "label       : %s\n", label)
	}
	fmt.Fprintf(w, "volume  V   : %s cm^3\n", fmtNum(r.Volume))
	fmt.Fprintf(w, "area    A   : %s cm^2\n", fmtNum(r.Area))
	fmt.Fprintf(w, "min area    : %s cm^2 (enclosing sphere)\n", fmtNum(r.MinArea))
	fmt.Fprintf(w, "modulus M   : %s cm   (M = V/A)\n", fmtMod(r.Modulus))
	fmt.Fprintf(w, "mold const C: %s min/cm^n\n", fmtNum(r.MoldConst))
	fmt.Fprintf(w, "exponent  n : %s\n", fmtNum(r.Exponent))
	fmt.Fprintf(w, "freeze time : %s min  (tf = C * M^n)\n", fmtNum(r.FreezeTime))
	if r.SuperheatApplied {
		fmt.Fprintf(w, "superheat   : +%s K applied\n", fmtNum(r.SuperheatK))
	}
}

func writeShapeTable(w io.Writer, t chvorinov.ShapeTable) {
	writeHeading(w, "shape comparison (same volume)")
	fmt.Fprintf(w, "%-10s %-18s %10s %10s %10s %10s\n",
		"shape", "size", "A cm^2", "M cm", "tf min", "ideal M")
	for _, row := range t.Rows {
		fmt.Fprintf(w, "%-10s %-18s %10s %10s %10s %10s  (M=%s)\n",
			row.Label,
			row.SizeDesc,
			fmtNum(row.Area),
			fmtMod(row.Modulus),
			fmtNum(row.Time),
			fmtMod(row.IdealModulus),
			row.IdealFormula,
		)
	}
	fmt.Fprintf(w, "slowest: %s   fastest: %s\n", t.Slowest, t.Fastest)
}

func writeRiser(w io.Writer, check chvorinov.RiserCheck) {
	writeHeading(w, "riser")
	fmt.Fprintf(w, "casting modulus : %s cm\n", fmtMod(check.CastingModulus))
	fmt.Fprintf(w, "riser modulus   : %s cm\n", fmtMod(check.RiserModulus))
	fmt.Fprintf(w, "required modulus: %s cm\n", fmtMod(check.RequiredModulus))
	fmt.Fprintf(w, "status          : %s\n", check.Status)
	if check.HasWarning() {
		fmt.Fprintf(w, "warning         : %s\n", check.Warning)
	}
}

func writeScale(w io.Writer, rep chvorinov.CrossScaleReport) {
	writeHeading(w, "similarity scaling")
	fmt.Fprintf(w, "factor      : %s (every linear dimension x2)\n", fmtNum(rep.Factor))
	fmt.Fprintf(w, "exponent    : %s\n", fmtNum(rep.Exponent))
	fmt.Fprintf(w, "modulus     : %s cm -> %s cm  (M x %s)\n",
		fmtMod(rep.RefModulus), fmtMod(rep.ScaledModulus), fmtNum(rep.Factor))
	fmt.Fprintf(w, "freeze time : %s min -> %s min  (ratio %s)\n",
		fmtNum(rep.RefTime), fmtNum(rep.ScaledTime), fmtNum(rep.TimeRatio))
	fmt.Fprintf(w, "cross rule  : size x2 with n=2 => tf x4 (here tf x %s)\n",
		fmtNum(rep.TimeRatio))
}

func writeVersion(w io.Writer, version string) {
	fmt.Fprintf(w, "chvorinov-t %s\n", version)
}

func writeWarning(w io.Writer, msg string) {
	fmt.Fprintf(w, "chvorinov-t: warning: %s\n", msg)
}
