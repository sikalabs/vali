package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var valiBin string

func TestMain(m *testing.M) {
	bin, err := os.MkdirTemp("", "vali-e2e-*")
	if err != nil {
		panic(err)
	}
	valiBin = filepath.Join(bin, "vali")
	if err := exec.Command("go", "build", "-o", valiBin, "github.com/sikalabs/vali").Run(); err != nil {
		panic("failed to build vali: " + err.Error())
	}
	code := m.Run()
	os.RemoveAll(bin)
	os.Exit(code)
}

func run(t *testing.T, env []string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(valiBin, args...)
	cmd.Env = append(os.Environ(), env...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			exitCode = exit.ExitCode()
		} else {
			t.Fatalf("unexpected exec error: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func writeConfig(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

// --- JSON config ---

func Test_Validate_JSON_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.json", `{"env":["MY_VAR"]}`)

	stdout, _, code := run(t, []string{"MY_VAR=hello"}, "validate", "--config", path)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "OK") {
		t.Errorf("expected OK in stdout, got %q", stdout)
	}
}

func Test_Validate_JSON_MissingEnv(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.json", `{"env":["MISSING_VAR"]}`)

	_, stderr, code := run(t, nil, "validate", "--config", path)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "MISSING_VAR") {
		t.Errorf("expected MISSING_VAR in stderr, got %q", stderr)
	}
}

func Test_Validate_JSON_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.json", `{"files":["/nonexistent/path"]}`)

	_, stderr, code := run(t, nil, "validate", "--config", path)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "/nonexistent/path") {
		t.Errorf("expected missing path in stderr, got %q", stderr)
	}
}

// --- YAML config ---

func Test_Validate_YAML_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.yaml", "env:\n  - MY_VAR\n")

	stdout, _, code := run(t, []string{"MY_VAR=hello"}, "validate", "--config", path)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "OK") {
		t.Errorf("expected OK in stdout, got %q", stdout)
	}
}

func Test_Validate_YAML_MissingEnv(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.yaml", "env:\n  - MISSING_VAR\n")

	_, stderr, code := run(t, nil, "validate", "--config", path)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "MISSING_VAR") {
		t.Errorf("expected MISSING_VAR in stderr, got %q", stderr)
	}
}

func Test_Validate_YAML_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "vali.yaml", "files:\n  - /nonexistent/path\n")

	_, stderr, code := run(t, nil, "validate", "--config", path)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "/nonexistent/path") {
		t.Errorf("expected missing path in stderr, got %q", stderr)
	}
}

// --- Auto-discovery ---

func Test_Validate_AutoDiscover_YAML(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	writeConfig(t, dir, "vali.yaml", "env:\n  - AUTO_VAR\n")

	stdout, _, code := run(t, []string{"AUTO_VAR=yes"}, "validate")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "OK") {
		t.Errorf("expected OK, got %q", stdout)
	}
}

func Test_Validate_AutoDiscover_LocalOverridesYAML(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	writeConfig(t, dir, "vali.yaml", "env:\n  - SHOULD_NOT_BE_CHECKED\n")
	writeConfig(t, dir, "vali.local.yaml", "env:\n  - LOCAL_VAR\n")

	_, stderr, code := run(t, nil, "validate")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "LOCAL_VAR") {
		t.Errorf("expected LOCAL_VAR in stderr (local config used), got %q", stderr)
	}
}

func Test_Validate_NoConfig(t *testing.T) {
	t.Chdir(t.TempDir())
	_, stderr, code := run(t, nil, "validate")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr, "no vali config file found") {
		t.Errorf("expected 'no vali config file found' in stderr, got %q", stderr)
	}
}
