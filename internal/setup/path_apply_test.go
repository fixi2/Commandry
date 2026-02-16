package setup

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEnsureWindowsUserPathConfiguredAddsEntry(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	prevRead := readWindowsUserPathFn
	prevWrite := writeWindowsUserPathFn
	defer func() {
		readWindowsUserPathFn = prevRead
		writeWindowsUserPathFn = prevWrite
	}()

	binDir := `C:\Users\test\AppData\Local\InfraTrack\bin`
	readWindowsUserPathFn = func() (string, error) {
		return `C:\Tools`, nil
	}
	var written string
	writeWindowsUserPathFn = func(v string) error {
		written = v
		return nil
	}

	res, err := ensureWindowsUserPathConfigured(binDir)
	if err != nil {
		t.Fatalf("ensureWindowsUserPathConfigured failed: %v", err)
	}
	if !res.Changed {
		t.Fatalf("expected changed=true")
	}
	if !PathContainsDir(written, binDir) {
		t.Fatalf("expected written PATH to contain %q, got %q", binDir, written)
	}
	if res.PathEntry != binDir {
		t.Fatalf("expected path entry %q, got %q", binDir, res.PathEntry)
	}
}

func TestEnsureWindowsUserPathConfiguredNoChangeWhenPresent(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	prevRead := readWindowsUserPathFn
	prevWrite := writeWindowsUserPathFn
	defer func() {
		readWindowsUserPathFn = prevRead
		writeWindowsUserPathFn = prevWrite
	}()

	binDir := `C:\Users\test\AppData\Local\InfraTrack\bin`
	readWindowsUserPathFn = func() (string, error) {
		return strings.Join([]string{`C:\Tools`, binDir}, ";"), nil
	}
	writeWindowsUserPathFn = func(v string) error {
		t.Fatalf("unexpected write: %s", v)
		return nil
	}

	res, err := ensureWindowsUserPathConfigured(binDir)
	if err != nil {
		t.Fatalf("ensureWindowsUserPathConfigured failed: %v", err)
	}
	if res.Changed {
		t.Fatalf("expected changed=false")
	}
}

func TestEnsurePosixUserPathConfiguredUsesMarkerBlock(t *testing.T) {
	prevResolve := resolvePosixProfileFn
	defer func() { resolvePosixProfileFn = prevResolve }()

	profile := filepath.Join(t.TempDir(), ".profile")
	resolvePosixProfileFn = func() (string, error) { return profile, nil }

	binDir := "/tmp/infratrack/bin"
	res, err := ensurePosixUserPathConfigured(binDir)
	if err != nil {
		t.Fatalf("ensurePosixUserPathConfigured failed: %v", err)
	}
	if !res.Changed {
		t.Fatalf("expected changed=true")
	}

	res2, err := ensurePosixUserPathConfigured(binDir)
	if err != nil {
		t.Fatalf("2nd ensurePosixUserPathConfigured failed: %v", err)
	}
	if res2.Changed {
		t.Fatalf("expected changed=false on second run")
	}
}
