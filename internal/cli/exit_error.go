package cli

import "fmt"

// ExitError keeps the child process exit code for CLI propagation.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return fmt.Sprintf("command failed with exit code %d", e.Code)
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

func (e *ExitError) ExitCode() int {
	if e.Code <= 0 {
		return 1
	}

	return e.Code
}
