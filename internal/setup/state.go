package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fixi2/Commandry/internal/appdir"
)

func DefaultStatePath() (string, error) {
	cfg, err := appdir.ResolveConfigRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, "setup-state.json"), nil
}

func LoadState(path string) (StateFile, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return StateFile{}, false, nil
		}
		return StateFile{}, false, fmt.Errorf("read setup state: %w", err)
	}

	var s StateFile
	if err := json.Unmarshal(data, &s); err != nil {
		return StateFile{}, false, fmt.Errorf("decode setup state: %w", err)
	}
	if s.SchemaVersion == 0 {
		s.SchemaVersion = StateSchemaVersion
	}
	return s, true, nil
}

func SaveState(path string, state StateFile) error {
	if state.SchemaVersion == 0 {
		state.SchemaVersion = StateSchemaVersion
	}
	if state.Timestamp.IsZero() {
		state.Timestamp = time.Now().UTC()
	}

	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode setup state: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create setup state dir: %w", err)
	}
	return os.WriteFile(path, payload, 0o600)
}
