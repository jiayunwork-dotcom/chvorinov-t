// This file defines the report formatting shared by the subcommands. All
// reports are plain text tables written to an io.Writer so they can be
// tested without spawning processes. Numbers are rounded to a fixed number
// of decimals so the output is stable across runs.
package cli

import (
	"fmt"
	"io"
	"strings"

	"chvorinov-t/internal/chvorinov"
)

// fmtNum renders a float with three decimals, trimming a trailing zero so
// the output stays compact for whole numbers.
func fmtNum(v float64) string {
	return fmt.Sprintf("%.3f", v)
}

// fmtMod renders a modulus with four decimals because moduli are often
// small compared with the volume and area.
func fmtMod(v float64) string {
	return fmt.Sprintf("%.4f", v)
}

// writeHeading prints a section heading with a full-width rule under it.
func writeHeading(w io.Writer, title string) {
	fmt.Fprintf(w, "\n%s\n%s\n", title, strings.Repeat("-", len(title)))
}

// writeResult prints the core quantities of a casting: geometry, modulus,
// mold constant, exponent, and freezing time. The shape comparison table
// is printed by a separate call so the compare subcommand can reuse it.
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

// writeShapeTable prints the per-shape comparison for a fixed volume.
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

// writeRiser prints the riser check and any warning that came out of it.
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

// writeScale prints the similarity-scaling cross rule check.
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

// writeVersion prints the tool name and version line.
func writeVersion(w io.Writer, version string) {
	fmt.Fprintf(w, "chvorinov-t %s\n", version)
}

// writeWarning prints a warning line to a writer (normally stderr). It is
// separate from writeError because a warning does not change the exit code.
func writeWarning(w io.Writer, msg string) {
	fmt.Fprintf(w, "chvorinov-t: warning: %s\n", msg)
}
