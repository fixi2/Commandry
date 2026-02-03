package store

import "time"

type Step struct {
	Timestamp  time.Time `json:"timestamp"`
	Command    string    `json:"command"`
	ExitCode   int       `json:"exit_code"`
	DurationMS int64     `json:"duration_ms"`
	CWD        string    `json:"cwd,omitempty"`
}

type Session struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Steps     []Step     `json:"steps"`
}
