package capture

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

type RunResult struct {
	StartedAt time.Time
	Duration  time.Duration
	ExitCode  int
}

// RunCommand executes the provided command without capturing stdout or stderr.
func RunCommand(ctx context.Context, args []string, cwd string) (RunResult, error) {
	startedAt := time.Now().UTC()
	result := RunResult{
		StartedAt: startedAt,
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	result.Duration = time.Since(startedAt)
	result.ExitCode = 0

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, err
	}

	result.ExitCode = 127
	return result, err
}
