package setup

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestUndoNoState(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)

	result, err := Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if result.Changed {
		t.Fatalf("expected no changes when state is missing")
	}
}

func TestUndoRemovesInstalledBinaryAndState(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)

	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	binPath := ResolveTargetBinaryPath(binDir)
	if err := os.WriteFile(binPath, []byte("x"), 0o700); err != nil {
		t.Fatalf("write bin: %v", err)
	}

	statePath, err := DefaultStatePath()
	if err != nil {
		t.Fatalf("DefaultStatePath failed: %v", err)
	}
	err = SaveState(statePath, StateFile{
		SchemaVersion:    StateSchemaVersion,
		InstalledBinPath: binPath,
		Timestamp:        time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	result, err := Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if !result.Changed {
		t.Fatalf("expected changed=true")
	}
	if _, err := os.Stat(binPath); !os.IsNotExist(err) {
		t.Fatalf("expected installed binary removed, err=%v", err)
	}
	if _, found, err := LoadState(statePath); err != nil || found {
		t.Fatalf("expected state removed (found=%v, err=%v)", found, err)
	}
}

func TestUndoRemovesPosixMarkerBlockFromTouchedFile(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)

	profile := filepath.Join(root, ".profile")
	before := strings.Join([]string{
		"export A=1",
		setupPathBeginMarker,
		"export PATH=\"/tmp/bin:$PATH\"",
		setupPathEndMarker,
		"export B=2",
		"",
	}, "\n")
	if err := os.WriteFile(profile, []byte(before), 0o600); err != nil {
		t.Fatalf("write profile: %v", err)
	}

	statePath, err := DefaultStatePath()
	if err != nil {
		t.Fatalf("DefaultStatePath failed: %v", err)
	}
	err = SaveState(statePath, StateFile{
		SchemaVersion:  StateSchemaVersion,
		FilesTouched:   []TouchedFile{{Path: profile, Marker: setupPathBeginMarker}},
		PathEntryAdded: "/tmp/bin",
		Timestamp:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	result, err := Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if !result.Changed {
		t.Fatalf("expected changed=true")
	}
	afterBytes, err := os.ReadFile(profile)
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	after := string(afterBytes)
	if strings.Contains(after, setupPathBeginMarker) || strings.Contains(after, setupPathEndMarker) {
		t.Fatalf("expected marker block removed, got: %q", after)
	}
}

func TestUndoWindowsPathEntryRemoval(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)

	prevRead := readWindowsUserPathFn
	prevWrite := writeWindowsUserPathFn
	defer func() {
		readWindowsUserPathFn = prevRead
		writeWindowsUserPathFn = prevWrite
	}()

	target := `C:\Users\test\AppData\Local\InfraTrack\bin`
	readWindowsUserPathFn = func() (string, error) {
		return strings.Join([]string{`C:\Tools`, target, `D:\Work`}, ";"), nil
	}
	var written string
	writeWindowsUserPathFn = func(v string) error {
		written = v
		return nil
	}

	statePath, err := DefaultStatePath()
	if err != nil {
		t.Fatalf("DefaultStatePath failed: %v", err)
	}
	err = SaveState(statePath, StateFile{
		SchemaVersion:  StateSchemaVersion,
		PathEntryAdded: target,
		Timestamp:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	result, err := Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if !result.Changed {
		t.Fatalf("expected changed=true")
	}
	if PathContainsDir(written, target) {
		t.Fatalf("expected target removed from written PATH: %q", written)
	}
}
