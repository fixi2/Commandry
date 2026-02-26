package appdir

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	CurrentDirName = "commandry"
	LegacyDirName  = "infratrack"
)

func ResolveConfigRoot() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}

	current := filepath.Join(base, CurrentDirName)
	legacy := filepath.Join(base, LegacyDirName)
	if _, err := os.Stat(current); err == nil {
		return current, nil
	}
	if _, err := os.Stat(legacy); err == nil {
		if err := MigrateLegacyDir(legacy, current); err == nil {
			return current, nil
		}
		return legacy, nil
	}
	return current, nil
}

func MigrateLegacyDir(legacyDir, currentDir string) error {
	if samePath(legacyDir, currentDir) {
		return nil
	}
	if _, err := os.Stat(currentDir); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat current config dir: %w", err)
	}
	if _, err := os.Stat(legacyDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat legacy config dir: %w", err)
	}

	if err := os.Rename(legacyDir, currentDir); err != nil {
		return fmt.Errorf("migrate legacy config dir: %w", err)
	}
	return nil
}

func samePath(a, b string) bool {
	aa := filepath.Clean(a)
	bb := filepath.Clean(b)
	return aa == bb
}
