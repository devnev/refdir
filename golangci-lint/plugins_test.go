package golangcilint

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPlugin(t *testing.T) {
	d := t.TempDir()

	cmd := exec.Command(`sh`, `-c`, `golangci-lint custom && mv custom-gcl "${D}"`)
	cmd.Dir = "./testdata"
	cmd.Env = append(os.Environ(), "D="+d)
	t.Logf("Building custom linter with %q", cmd.Args)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build custom linter: %v, output: %s", err, string(out))
	}
	binPath := filepath.Join(d, "custom-gcl")

	t.Logf("Cleaning cache")
	if out, err := exec.Command(binPath, "cache", "clean").CombinedOutput(); err != nil {
		t.Fatalf("Failed to clear lint cache: %v, output: %s", err, string(out))
	}

	// Should succeed

	t.Logf("Running pass checks")

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-default.yml", "../analysis/refdir/internal/example_default_good")
	if _, err := cmd.Output(); err != nil {
		t.Errorf("Custom linter failed on default example with default config: %s", fmtExitError(err))
	}

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-funcup.yml", "../analysis/refdir/internal/example_funcup_good")
	if _, err := cmd.Output(); err != nil {
		t.Errorf("Expected custom linter to fail on funcup example with funcup config, got %s", fmtExitError(err))
	}

	// Should fail

	t.Logf("Running fail checks")

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-default.yml", "../analysis/refdir/internal/example_funcup_good")
	if _, err := cmd.Output(); exitCode(err) != 1 {
		t.Errorf("Expected custom linter to fail on funcup example with default config, got %s", fmtExitError(err))
	}

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-default.yml", "../analysis/refdir/internal/example_bad")
	if _, err := cmd.Output(); exitCode(err) != 1 {
		t.Errorf("Expected custom linter to fail on bad example with default config, got %s", fmtExitError(err))
	}

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-funcup.yml", "../analysis/refdir/internal/example_default_good")
	if _, err := cmd.Output(); exitCode(err) != 1 {
		t.Errorf("Custom linter failed on default example with funcup config: %s", fmtExitError(err))
	}

	cmd = exec.Command(binPath, "run", "--config=testdata/.golangci-funcup.yml", "../analysis/refdir/internal/example_bad")
	if _, err := cmd.Output(); exitCode(err) != 1 {
		t.Errorf("Expected custom linter to fail on bad example with default config, got %s", fmtExitError(err))
	}
}

func exitCode(err error) int {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return -1
	}
	return exitErr.ExitCode()
}

func fmtExitError(err error) string {
	if err == nil {
		return fmt.Sprint(nil)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || len(exitErr.Stderr) == 0 {
		return err.Error()
	}
	return fmt.Sprintf("%v (stderr: %s)", err, string(exitErr.Stderr))
}
