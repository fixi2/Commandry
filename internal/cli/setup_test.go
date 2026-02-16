package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestSetupCommandsExist(t *testing.T) {
	t.Parallel()

	root, err := NewRootCommand()
	if err != nil {
		t.Fatalf("NewRootCommand failed: %v", err)
	}

	setupCmd, _, err := root.Find([]string{"setup"})
	if err != nil {
		t.Fatalf("root.Find(setup) failed: %v", err)
	}
	if setupCmd == nil || setupCmd.Name() != "setup" {
		t.Fatalf("setup command not found")
	}

	setupStatus, _, err := root.Find([]string{"setup", "status"})
	if err != nil {
		t.Fatalf("root.Find(setup status) failed: %v", err)
	}
	if setupStatus == nil || setupStatus.Name() != "status" {
		t.Fatalf("setup status command not found")
	}
}

func TestSetupDryRunOutput(t *testing.T) {
	t.Parallel()

	root, err := NewRootCommand()
	if err != nil {
		t.Fatalf("NewRootCommand failed: %v", err)
	}
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"setup", "--scope", "user", "--completion", "none"})

	if err := root.Execute(); err != nil {
		t.Fatalf("setup dry-run failed: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Dry-run only") {
		t.Fatalf("expected dry-run text, got: %s", got)
	}
	if !strings.Contains(got, "Planned actions:") {
		t.Fatalf("expected actions section, got: %s", got)
	}
}

func TestSetupStatusJSON(t *testing.T) {
	t.Parallel()

	root, err := NewRootCommand()
	if err != nil {
		t.Fatalf("NewRootCommand failed: %v", err)
	}
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"setup", "status", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("setup status --json failed: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("status json decode failed: %v\npayload:\n%s", err, out.String())
	}
	if _, ok := decoded["installed"]; !ok {
		t.Fatalf("expected installed field in json: %v", decoded)
	}
	if _, ok := decoded["pathOk"]; !ok {
		t.Fatalf("expected pathOk field in json: %v", decoded)
	}
}
