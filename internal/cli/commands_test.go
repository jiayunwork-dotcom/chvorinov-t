package cli

import (
	"bytes"
	"strings"
	"testing"
)

// runCLI invokes Run with captured writers and returns the exit code, the
// stdout text, and the stderr text.
func runCLI(args ...string) (code int, stdout, stderr string) {
	var out, errBuf bytes.Buffer
	code = Run(args, &out, &errBuf)
	return code, out.String(), errBuf.String()
}

// TestRunFreezePrintsModulusAndTime checks the main command path: a valid
// steel cube prints a modulus and a freezing time and exits 0.
func TestRunFreezePrintsModulusAndTime(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"label":"steel-cube","shape":"cube","edge":10,"C":3,"n":2}`)
	code, stdout, stderr := runCLI("freeze", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "freeze time") {
		t.Errorf("stdout missing freeze time:\n%s", stdout)
	}
	if !strings.Contains(stdout, "modulus") {
		t.Errorf("stdout missing modulus:\n%s", stdout)
	}
	if !strings.Contains(stdout, "shape comparison") {
		t.Errorf("stdout missing shape comparison table:\n%s", stdout)
	}
}

// TestRunFreezeCubeModulusIsEdgeOverSix verifies the printed modulus for
// a 10 cm cube is 10/6, the a/6 reference.
func TestRunFreezeCubeModulusIsEdgeOverSix(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3,"n":2}`)
	code, stdout, _ := runCLI("freeze", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "1.6667") {
		t.Errorf("stdout should contain modulus 1.6667:\n%s", stdout)
	}
}

// TestRunFreezeZeroVolumeFails pins the "volume = 0" failure case: the
// command must exit non-zero and write an error to stderr.
func TestRunFreezeZeroVolumeFails(t *testing.T) {
	p := writeTempJSON(t, "zero.json", `{"shape":"custom","volume":0,"area":600,"C":3,"n":2}`)
	code, _, stderr := runCLI("freeze", p)
	if code == ExitOK {
		t.Fatal("exit = 0, want non-zero for zero volume")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunFreezeNegativeAreaFails pins the negative-area failure case.
func TestRunFreezeNegativeAreaFails(t *testing.T) {
	p := writeTempJSON(t, "neg.json", `{"shape":"custom","volume":1000,"area":-5,"C":3,"n":2}`)
	code, _, stderr := runCLI("freeze", p)
	if code == ExitOK {
		t.Fatal("exit = 0, want non-zero for negative area")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunFreezeZeroMoldConstFails pins the zero-C failure case.
func TestRunFreezeZeroMoldConstFails(t *testing.T) {
	p := writeTempJSON(t, "zeroC.json", `{"shape":"cube","edge":10,"C":0,"n":2}`)
	code, _, stderr := runCLI("freeze", p)
	if code == ExitOK {
		t.Fatal("exit = 0, want non-zero for C = 0")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunMissingFileFails pins the file-not-found failure case.
func TestRunMissingFileFails(t *testing.T) {
	code, _, stderr := runCLI("freeze", "/nonexistent/cube.json")
	if code == ExitOK {
		t.Fatal("exit = 0, want non-zero for a missing file")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunUnknownCommandIsUsageError pins the unknown-subcommand case as a
// usage error (exit 2).
func TestRunUnknownCommandIsUsageError(t *testing.T) {
	code, _, stderr := runCLI("pour")
	if code != ExitUsage {
		t.Errorf("exit = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("stderr has no unknown-command message: %q", stderr)
	}
}

// TestRunFreezeMissingArgumentIsUsageError pins the no-file case.
func TestRunFreezeMissingArgumentIsUsageError(t *testing.T) {
	code, _, stderr := runCLI("freeze")
	if code != ExitUsage {
		t.Errorf("exit = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr, "JSON input file") {
		t.Errorf("stderr has no usage hint: %q", stderr)
	}
}

// TestRunVersionOK pins the version subcommand.
func TestRunVersionOK(t *testing.T) {
	code, stdout, _ := runCLI("version")
	if code != ExitOK {
		t.Errorf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "chvorinov-t") {
		t.Errorf("stdout = %q, want version banner", stdout)
	}
}

// TestRunHelpOK pins --help as a zero-exit usage printer.
func TestRunHelpOK(t *testing.T) {
	code, stdout, _ := runCLI("--help")
	if code != ExitOK {
		t.Errorf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "freeze") {
		t.Errorf("help output missing subcommand list: %q", stdout)
	}
}

// TestRunRiserWarningGoesToStderr pins the undersized-riser warning: the
// command still succeeds (exit 0) but the warning is visible on stderr.
func TestRunRiserWarningGoesToStderr(t *testing.T) {
	p := writeTempJSON(t, "riser.json", `{
		"shape":"cube","edge":10,"C":3,"n":2,
		"riser":{"shape":"cylinder","radius":2,"height":4}
	}`)
	code, _, stderr := runCLI("freeze", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0 (warning must not fail the command)", code)
	}
	if !strings.Contains(stderr, "warning") {
		t.Errorf("stderr has no warning line: %q", stderr)
	}
}

// TestRunRiserSubcommandOK pins the standalone riser command.
func TestRunRiserSubcommandOK(t *testing.T) {
	p := writeTempJSON(t, "riser.json", `{
		"shape":"cube","edge":10,"C":3,
		"riser":{"shape":"cylinder","radius":6,"height":30}
	}`)
	code, stdout, _ := runCLI("riser", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "status") {
		t.Errorf("riser output missing status line:\n%s", stdout)
	}
}

// TestRunRiserWithoutRiserFails pins the missing-riser case for the riser
// subcommand.
func TestRunRiserWithoutRiserFails(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3}`)
	code, _, stderr := runCLI("riser", p)
	if code == ExitOK {
		t.Fatal("exit = 0, want non-zero for a file without a riser")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunScaleOK pins the scale subcommand output.
func TestRunScaleOK(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3,"n":2}`)
	code, stdout, _ := runCLI("scale", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "similarity scaling") {
		t.Errorf("scale output missing heading:\n%s", stdout)
	}
}

// TestRunValidateOKAndFails pins validate's success and failure paths.
func TestRunValidateOKAndFails(t *testing.T) {
	ok := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3}`)
	code, stdout, _ := runCLI("validate", ok)
	if code != ExitOK {
		t.Fatalf("validate valid file exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "valid") {
		t.Errorf("validate output = %q, want a valid line", stdout)
	}

	bad := writeTempJSON(t, "bad.json", `{"shape":"cube","edge":0,"C":3}`)
	code, _, stderr := runCLI("validate", bad)
	if code == ExitOK {
		t.Fatal("validate invalid file exit = 0, want non-zero")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr has no error line: %q", stderr)
	}
}

// TestRunCompareOK pins the compare subcommand output.
func TestRunCompareOK(t *testing.T) {
	p := writeTempJSON(t, "cube.json", `{"shape":"cube","edge":10,"C":3,"n":2}`)
	code, stdout, _ := runCLI("compare", p)
	if code != ExitOK {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "shape comparison") {
		t.Errorf("compare output missing heading:\n%s", stdout)
	}
}
