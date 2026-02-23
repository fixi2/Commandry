package setup

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestApplyInstallsBinaryAndWritesState(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", filepath.Join(root, "AppData", "Roaming"))

	source := filepath.Join(root, "source-bin")
	if runtime.GOOS == "windows" {
		source += ".exe"
	}
	if err := os.WriteFile(source, []byte("infratrack-binary"), 0o700); err != nil {
		t.Fatalf("write source failed: %v", err)
	}

	binDir := filepath.Join(root, "bin")
	result, err := Apply(ApplyInput{
		Scope:            ScopeUser,
		BinDir:           binDir,
		NoPath:           true,
		Completion:       CompletionNone,
		SourceBinaryPath: source,
	})
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if result.InstalledBinPath == "" {
		t.Fatalf("expected installed binary path")
	}
	if _, err := os.Stat(result.InstalledBinPath); err != nil {
		t.Fatalf("installed binary not found: %v", err)
	}
	if _, found, err := LoadState(result.StatePath); err != nil || !found {
		t.Fatalf("expected state file (found=%v, err=%v)", found, err)
	}
}

func TestApplyWindowsStagingName(t *testing.T) {
	target := `C:\Users\me\AppData\Local\InfraTrack\bin\infratrack.exe`
	got := windowsStagingPath(target)
	if !strings.HasSuffix(strings.ToLower(got), "infratrack.new.exe") {
		t.Fatalf("unexpected staging path: %s", got)
	}
}

func TestApplyRejectsControlCharsInBinDir(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", filepath.Join(root, "AppData", "Roaming"))

	source := filepath.Join(root, "source-bin")
	if runtime.GOOS == "windows" {
		source += ".exe"
	}
	if err := os.WriteFile(source, []byte("infratrack-binary"), 0o700); err != nil {
		t.Fatalf("write source failed: %v", err)
	}

	_, err := Apply(ApplyInput{
		Scope:            ScopeUser,
		BinDir:           "/tmp/infratrack\nbin",
		NoPath:           true,
		Completion:       CompletionNone,
		SourceBinaryPath: source,
	})
	if err == nil {
		t.Fatal("expected apply to reject control characters in --bin-dir")
	}
}

func TestApplyWindowsPreservesPathOwnershipAcrossIdempotentApply(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", filepath.Join(root, "AppData", "Roaming"))

	prevRead := readWindowsUserPathFn
	prevWrite := writeWindowsUserPathFn
	defer func() {
		readWindowsUserPathFn = prevRead
		writeWindowsUserPathFn = prevWrite
	}()

	userPath := `C:\Tools`
	readWindowsUserPathFn = func() (string, error) { return userPath, nil }
	writeWindowsUserPathFn = func(v string) error {
		userPath = v
		return nil
	}

	source := filepath.Join(root, "source-bin.exe")
	if err := os.WriteFile(source, []byte("infratrack-binary"), 0o700); err != nil {
		t.Fatalf("write source failed: %v", err)
	}
	binDir := filepath.Join(root, "bin")

	first, err := Apply(ApplyInput{
		Scope:            ScopeUser,
		BinDir:           binDir,
		NoPath:           false,
		Completion:       CompletionNone,
		SourceBinaryPath: source,
	})
	if err != nil {
		t.Fatalf("1st Apply failed: %v", err)
	}
	if first.PathEntryAdded == "" {
		t.Fatalf("expected path entry to be recorded on first apply")
	}

	second, err := Apply(ApplyInput{
		Scope:            ScopeUser,
		BinDir:           binDir,
		NoPath:           false,
		Completion:       CompletionNone,
		SourceBinaryPath: source,
	})
	if err != nil {
		t.Fatalf("2nd Apply failed: %v", err)
	}
	if second.PathEntryAdded == "" {
		t.Fatalf("expected path ownership to persist on second apply")
	}

	state, found, err := LoadState(second.StatePath)
	if err != nil || !found {
		t.Fatalf("expected state after second apply (found=%v, err=%v)", found, err)
	}
	if normalizePathForCompare(state.PathEntryAdded) != normalizePathForCompare(binDir) {
		t.Fatalf("expected state path entry to match bin dir, got %q", state.PathEntryAdded)
	}

	_, err = Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if PathContainsDir(userPath, binDir) {
		t.Fatalf("expected undo to remove bin dir from user PATH, got %q", userPath)
	}
}

func TestFinalizeWindowsBinaryRollsBackOnActivationFailure(t *testing.T) {
	target := `C:\bin\infratrack.exe`
	staging := `C:\bin\infratrack.new.exe`
	backup := target + ".bak"

	origRename := osRenameFn
	origRemove := osRemoveFn
	defer func() {
		osRenameFn = origRename
		osRemoveFn = origRemove
	}()

	osRemoveFn = func(string) error { return nil }
	renamedToBackup := false
	rolledBack := false
	osRenameFn = func(oldPath, newPath string) error {
		if oldPath == target && newPath == backup {
			renamedToBackup = true
			return nil
		}
		if oldPath == staging && newPath == target {
			return os.ErrNotExist
		}
		if oldPath == backup && newPath == target {
			rolledBack = true
			return nil
		}
		return nil
	}

	_, _, err := finalizeWindowsBinary(staging, target)
	if err == nil {
		t.Fatal("expected finalizeWindowsBinary to fail")
	}
	if !renamedToBackup {
		t.Fatal("expected existing target to be moved to backup before activation")
	}
	if !rolledBack {
		t.Fatal("expected rollback rename from backup to target")
	}
}
