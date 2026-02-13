package cli

import (
	"strings"
	"testing"

	"github.com/fixi2/InfraTrack/internal/store"
)

func TestCollectFlaggedSteps(t *testing.T) {
	t.Parallel()

	steps := []store.Step{
		{Command: "ok command", Status: "OK", ExitCode: intPtr(0)},
		{Command: "[REDACTED BY POLICY]", Status: "REDACTED", Reason: "policy_redacted"},
		{Command: "failed command", Status: "FAILED", ExitCode: intPtr(1)},
		{Command: "nonzero no status", ExitCode: intPtr(2)},
	}

	session := &store.Session{Steps: steps}
	flagged := collectFlaggedSteps(session)
	if len(flagged) != 3 {
		t.Fatalf("expected 3 flagged steps, got %d", len(flagged))
	}
	if flagged[0].Number != 1 || flagged[1].Number != 2 || flagged[2].Number != 3 {
		t.Fatalf("unexpected numbering: %#v", flagged)
	}
	if !strings.Contains(flagged[0].Command, "[REDACTED BY POLICY]") {
		t.Fatalf("expected redacted command preview, got %q", flagged[0].Command)
	}
}

func TestParseSelection(t *testing.T) {
	t.Parallel()

	flagged := []flaggedStep{
		{Number: 1, StepIndex: 3},
		{Number: 2, StepIndex: 5},
	}

	global, steps, err := parseSelection("0 2 2 1", flagged)
	if err != nil {
		t.Fatalf("parseSelection returned error: %v", err)
	}
	if !global {
		t.Fatalf("expected global selection")
	}
	if len(steps) != 2 || steps[0] != 3 || steps[1] != 5 {
		t.Fatalf("unexpected step indexes: %#v", steps)
	}

	if _, _, err := parseSelection("9", flagged); err == nil {
		t.Fatalf("expected invalid selection error")
	}
}

func TestSanitizeComment(t *testing.T) {
	t.Parallel()

	raw := "  hello\r\nworld\t\x1b[31m \x00ok  "
	got := sanitizeComment(raw)
	if strings.Contains(got, "\n") || strings.Contains(got, "\r") {
		t.Fatalf("expected single line, got %q", got)
	}
	if strings.Contains(got, "\x1b") || strings.Contains(got, "\x00") {
		t.Fatalf("control chars were not removed: %q", got)
	}
	if got != "hello world [31m ok" {
		t.Fatalf("unexpected sanitized value: %q", got)
	}

	long := strings.Repeat("a", maxCommentLen+10)
	if len([]rune(sanitizeComment(long))) != maxCommentLen {
		t.Fatalf("expected truncated comment length %d", maxCommentLen)
	}
}

func intPtr(v int) *int { return &v }
