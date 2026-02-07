package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDoctorCommandExists(t *testing.T) {
	t.Parallel()

	root, err := NewRootCommand()
	if err != nil {
		t.Fatalf("NewRootCommand failed: %v", err)
	}

	cmd, _, err := root.Find([]string{"doctor"})
	if err != nil {
		t.Fatalf("root.Find(doctor) failed: %v", err)
	}
	if cmd == nil || cmd.Name() != "doctor" {
		t.Fatalf("doctor command not found")
	}
}

func TestDoctorCommandOutput(t *testing.T) {
	t.Parallel()

	root, err := NewRootCommand()
	if err != nil {
		t.Fatalf("NewRootCommand failed: %v", err)
	}

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"doctor"})

	if err := root.Execute(); err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"InfraTrack doctor",
		"Root dir:",
		"Tool availability:",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("doctor output missing %q in %q", want, text)
		}
	}
}

func TestPathContainsDir(t *testing.T) {
	t.Parallel()

	sep := string(os.PathListSeparator)
	pathEnv := `C:\bin` + sep + `C:\Tools` + sep + `C:\Windows`

	if !pathContainsDir(pathEnv, `C:\Tools`) {
		t.Fatalf("expected pathContainsDir to find directory")
	}
	if pathContainsDir(pathEnv, `C:\Missing`) {
		t.Fatalf("did not expect pathContainsDir to find missing directory")
	}
}
