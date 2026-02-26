package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolveScope(v string) (Scope, error) {
	s := Scope(strings.ToLower(strings.TrimSpace(v)))
	switch s {
	case ScopeUser, ScopeSystem:
		return s, nil
	default:
		return "", fmt.Errorf("unsupported scope %q (use user)", v)
	}
}

func ResolveCompletion(v string) (CompletionMode, error) {
	mode := CompletionMode(strings.ToLower(strings.TrimSpace(v)))
	switch mode {
	case "", CompletionNone:
		return CompletionNone, nil
	default:
		return "", fmt.Errorf("unsupported completion mode %q", v)
	}
}

func DefaultBinDir() (string, error) {
	if runtime.GOOS == "windows" {
		if localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); localAppData != "" {
			return filepath.Join(localAppData, "Commandry", "bin"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home dir: %w", err)
		}
		return filepath.Join(home, "AppData", "Local", "Commandry", "bin"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".local", "bin"), nil
}

func CurrentExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}
	return filepath.Clean(exe), nil
}

func ResolveTargetBinaryPath(binDir string) string {
	name := "cmdry"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(binDir, name)
}

func PathContainsDir(pathEnv, target string) bool {
	if strings.TrimSpace(target) == "" {
		return false
	}
	normalizedTarget := normalizePathForCompare(target)
	for _, part := range filepath.SplitList(pathEnv) {
		if normalizePathForCompare(part) == normalizedTarget {
			return true
		}
	}
	return false
}

func normalizePathForCompare(v string) string {
	s := strings.TrimSpace(v)
	if s == "" {
		return ""
	}
	s = filepath.Clean(s)
	s = strings.TrimRight(s, `\/`)
	if runtime.GOOS == "windows" {
		s = strings.ReplaceAll(s, "/", `\`)
		s = strings.ToLower(s)
	}
	return s
}

func validatePathNoControlChars(v string) error {
	if strings.Contains(v, "\x00") || strings.Contains(v, "\n") || strings.Contains(v, "\r") {
		return fmt.Errorf("path contains unsupported control characters")
	}
	return nil
}
