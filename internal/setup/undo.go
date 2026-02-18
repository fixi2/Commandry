package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Undo() (UndoResult, error) {
	statePath, err := DefaultStatePath()
	if err != nil {
		return UndoResult{}, err
	}

	result := UndoResult{
		StatePath: statePath,
		Actions:   make([]string, 0, 8),
	}

	state, found, err := LoadState(statePath)
	if err != nil {
		return result, err
	}
	if !found {
		return result, nil
	}

	if state.InstalledBinPath != "" {
		removed, err := removeFileIfExists(state.InstalledBinPath)
		if err != nil {
			return result, fmt.Errorf("remove installed binary: %w", err)
		}
		if removed {
			result.Changed = true
			result.Actions = append(result.Actions, fmt.Sprintf("Removed %s.", state.InstalledBinPath))
		}
	}

	if state.PathEntryAdded != "" {
		changed, action, err := undoPathEntry(state.PathEntryAdded)
		if err != nil {
			return result, err
		}
		if changed {
			result.Changed = true
		}
		if action != "" {
			result.Actions = append(result.Actions, action)
		}
	}

	for _, touched := range state.FilesTouched {
		if strings.TrimSpace(touched.Path) == "" {
			continue
		}
		changed, err := removeMarkerBlock(touched.Path, setupPathBeginMarker, setupPathEndMarker)
		if err != nil {
			return result, err
		}
		if changed {
			result.Changed = true
			result.Actions = append(result.Actions, fmt.Sprintf("Removed setup marker block from %s.", touched.Path))
		}
	}

	for i := len(state.CreatedDirs) - 1; i >= 0; i-- {
		dir := strings.TrimSpace(state.CreatedDirs[i])
		if dir == "" {
			continue
		}
		removed, err := removeDirIfEmpty(dir)
		if err != nil {
			return result, err
		}
		if removed {
			result.Changed = true
			result.Actions = append(result.Actions, fmt.Sprintf("Removed empty directory %s.", dir))
		}
	}

	if err := os.Remove(statePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return result, fmt.Errorf("remove setup state: %w", err)
	}
	if found {
		result.Changed = true
		result.Actions = append(result.Actions, "Removed setup state file.")
	}

	return result, nil
}

func undoPathEntry(pathEntry string) (bool, string, error) {
	if runtime.GOOS == "windows" {
		return undoWindowsUserPathEntry(pathEntry)
	}
	return undoPosixPathEntry(pathEntry)
}

func undoWindowsUserPathEntry(pathEntry string) (bool, string, error) {
	current, err := readWindowsUserPathFn()
	if err != nil {
		return false, "", fmt.Errorf("read user PATH: %w", err)
	}
	parts := filepath.SplitList(current)
	nextParts := make([]string, 0, len(parts))
	removed := false
	target := normalizePathForCompare(pathEntry)
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if normalizePathForCompare(p) == target {
			removed = true
			continue
		}
		nextParts = append(nextParts, p)
	}
	if !removed {
		return false, "", nil
	}
	next := strings.Join(nextParts, string(os.PathListSeparator))
	if err := writeWindowsUserPathFn(next); err != nil {
		return false, "", fmt.Errorf("write user PATH: %w", err)
	}
	return true, fmt.Sprintf("Removed %s from user PATH.", pathEntry), nil
}

func undoPosixPathEntry(_ string) (bool, string, error) {
	// POSIX changes are reverted via marker block removal from filesTouched.
	return false, "", nil
}

func removeMarkerBlock(path, beginMarker, endMarker string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("read profile %s: %w", path, err)
	}
	text := string(content)
	begin := strings.Index(text, beginMarker)
	if begin < 0 {
		return false, nil
	}
	end := strings.Index(text[begin:], endMarker)
	if end < 0 {
		return false, nil
	}
	end = begin + end + len(endMarker)
	for end < len(text) && (text[end] == '\r' || text[end] == '\n') {
		end++
	}
	next := strings.TrimRight(text[:begin]+text[end:], "\r\n")
	if next != "" {
		next += "\n"
	}
	if err := os.WriteFile(path, []byte(next), 0o600); err != nil {
		return false, fmt.Errorf("write profile %s: %w", path, err)
	}
	return true, nil
}

func removeFileIfExists(path string) (bool, error) {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func removeDirIfEmpty(path string) (bool, error) {
	err := os.Remove(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if errors.Is(err, os.ErrPermission) || strings.Contains(strings.ToLower(err.Error()), "directory not empty") {
		return false, nil
	}
	return false, err
}
