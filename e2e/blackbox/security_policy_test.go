package blackbox

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Contract: C3
func TestSecretsAndDenylistDoNotLeakToRunbook(t *testing.T) {
	t.Parallel()
	h := newHarness(t)

	h.initSession("security-blackbox")
	secret := randomSentinel()

	echoArgs := shellEchoCommand("--token=" + secret)
	h.run(append([]string{"run", "--"}, echoArgs...)...)
	// Known denylist contract example.
	h.run("run", "--", "printenv")

	h.stopSession()
	runbookPath := h.exportLastMD()
	runbook := readFile(t, runbookPath)

	if strings.Contains(runbook, secret) {
		t.Fatalf("secret leaked into runbook")
	}
	if !strings.Contains(runbook, "[REDACTED]") && !strings.Contains(runbook, "[REDACTED BY POLICY]") {
		t.Fatalf("expected policy redaction markers in runbook")
	}
	if !strings.Contains(runbook, "[REDACTED BY POLICY]") {
		t.Fatalf("expected denylist redaction marker in runbook")
	}
}

// Contract: C3 (property-style)
func TestRandomSentinelNeverAppearsInArtifacts(t *testing.T) {
	t.Parallel()
	h := newHarness(t)
	h.initSession("property-security")

	sentinels := make([]string, 0, 25)
	for i := 0; i < 25; i++ {
		s := randomSentinel()
		sentinels = append(sentinels, s)
		echoArgs := shellEchoCommand("--password=" + s)
		h.run(append([]string{"run", "--"}, echoArgs...)...)
	}
	h.stopSession()
	runbookPath := h.exportLastMD()

	content := readFile(t, runbookPath)
	for _, s := range sentinels {
		if strings.Contains(content, s) {
			t.Fatalf("sentinel leaked: %s", s)
		}
	}

	// Also scan config/store area in isolated env.
	paths := []string{
		filepath.Join(h.rootDir, "appdata"),
		filepath.Join(h.workDir, "runbooks"),
	}
	for _, root := range paths {
		allText := readAllFiles(t, root)
		for _, s := range sentinels {
			if strings.Contains(allText, s) {
				t.Fatalf("sentinel leaked in artifacts under %s: %s", root, s)
			}
		}
	}
}

// Contract: C3
func TestURIUserinfoIsRedacted(t *testing.T) {
	t.Parallel()
	h := newHarness(t)

	h.initSession("uri-redaction")
	args := append([]string{"run", "--"}, shellEchoCommand("https://alice:secret-pass@example.com/api")...)
	h.run(args...)
	h.stopSession()
	runbookPath := h.exportLastMD()
	runbook := readFile(t, runbookPath)

	if strings.Contains(runbook, "secret-pass") || strings.Contains(runbook, "alice:") {
		t.Fatalf("uri userinfo leaked in runbook")
	}
	if !strings.Contains(runbook, "https://[REDACTED]:[REDACTED]@example.com/api") {
		t.Fatalf("expected uri userinfo redaction marker in runbook")
	}
}

func readAllFiles(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// Skip binaries.
		if runtime.GOOS == "windows" && strings.HasSuffix(strings.ToLower(path), ".exe") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		b.Write(data)
		b.WriteString("\n")
		return nil
	})
	return b.String()
}
