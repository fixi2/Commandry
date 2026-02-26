package appdir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateLegacyDirMovesDataToCurrent(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	legacy := filepath.Join(base, LegacyDirName)
	current := filepath.Join(base, CurrentDirName)

	if err := os.MkdirAll(legacy, 0o700); err != nil {
		t.Fatalf("mkdir legacy: %v", err)
	}
	cfg := filepath.Join(legacy, "config.yaml")
	if err := os.WriteFile(cfg, []byte("policy:\n"), 0o600); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	if err := MigrateLegacyDir(legacy, current); err != nil {
		t.Fatalf("MigrateLegacyDir failed: %v", err)
	}

	migratedCfg := filepath.Join(current, "config.yaml")
	data, err := os.ReadFile(migratedCfg)
	if err != nil {
		t.Fatalf("read migrated config: %v", err)
	}
	if string(data) != "policy:\n" {
		t.Fatalf("unexpected migrated content: %q", string(data))
	}
}

func TestResolveConfigRootUsesCommandryDir(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("APPDATA", base)

	root, err := ResolveConfigRoot()
	if err != nil {
		t.Fatalf("ResolveConfigRoot failed: %v", err)
	}
	want := filepath.Join(base, CurrentDirName)
	if filepath.Clean(root) != filepath.Clean(want) {
		t.Fatalf("root mismatch: got %q want %q", root, want)
	}
}
