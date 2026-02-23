package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseConfigUsesDefaultsWhenPolicySectionMissing(t *testing.T) {
	t.Parallel()

	cfg, err := ParseConfig("capture:\n  include_stdout: false\n")
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}
	if len(cfg.Denylist) == 0 || len(cfg.RedactionKeywords) == 0 {
		t.Fatalf("expected defaults when policy section missing")
	}
	if cfg.EnforceDenylist {
		t.Fatalf("expected enforce_denylist default false")
	}
}

func TestLoadFromConfigParsesRuntimePolicy(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "config.yaml")
	content := strings.Join([]string{
		"policy:",
		"  denylist:",
		"    - docker login",
		"  redaction_keywords:",
		"    - session_token",
		"  enforce_denylist: true",
	}, "\n")
	if err := osWriteFile(path, []byte(content)); err != nil {
		t.Fatalf("write config: %v", err)
	}

	p, err := LoadFromConfig(path)
	if err != nil {
		t.Fatalf("LoadFromConfig failed: %v", err)
	}
	if !p.EnforceDenylist() {
		t.Fatalf("expected enforce_denylist true")
	}

	denied := p.Apply("docker login registry.example.com", []string{"docker", "login", "registry.example.com"})
	if !denied.Denied || denied.Command != DeniedPlaceholder {
		t.Fatalf("expected denylist to deny docker login")
	}

	sanitized := p.Apply("run --session-token=abc123", []string{"run", "--session-token=abc123"})
	if strings.Contains(sanitized.Command, "abc123") {
		t.Fatalf("expected session token to be redacted, got %q", sanitized.Command)
	}
}

func TestParseConfigRejectsInvalidEnforceValue(t *testing.T) {
	t.Parallel()

	_, err := ParseConfig(strings.Join([]string{
		"policy:",
		"  enforce_denylist: maybe",
	}, "\n"))
	if err == nil {
		t.Fatalf("expected parse error for invalid enforce_denylist")
	}
}

func TestParseConfigHandlesUTF8BOM(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		"\ufeffpolicy:",
		"  denylist:",
		"    - \"echo blocked-now\"",
		"  redaction_keywords:",
		"    - token",
		"  enforce_denylist: true",
	}, "\n")
	cfg, err := ParseConfig(content)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}
	if !cfg.EnforceDenylist {
		t.Fatalf("expected enforce_denylist true from BOM-prefixed config")
	}
	if len(cfg.Denylist) != 1 || cfg.Denylist[0] != "echo blocked-now" {
		t.Fatalf("unexpected denylist parsed from BOM config: %+v", cfg.Denylist)
	}
}

func osWriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}
