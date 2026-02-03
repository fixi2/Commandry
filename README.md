# InfraTrack v0.1 (MVP)

InfraTrack is a local-first CLI that records explicitly executed commands and exports a deterministic markdown runbook.

## Build

Requirements:
- Go 1.22 or newer

Commands:

```bash
go mod tidy
go build ./cmd/infratrack
```

The binary is created as `infratrack` (or `infratrack.exe` on Windows).

## Quickstart Demo

```bash
infratrack init
infratrack start "Deploy to staging"
infratrack run -- kubectl apply -f deploy.yaml
infratrack run -- kubectl rollout status deploy/api
infratrack stop
infratrack export --last --md
```

Expected output artifact:
- `runbooks/<timestamp>-<slug>.md`

## CLI Commands

- `infratrack init` initializes local config and session storage in `os.UserConfigDir()/infratrack`.
- `infratrack start "<title>"` starts recording session metadata.
- `infratrack run -- <cmd ...>` executes command and records a sanitized step.
- `infratrack status` shows current recording state.
- `infratrack stop` finalizes the active session.
- `infratrack export --last --md` exports the latest completed session to markdown.

## Security Notes

- Recording is off by default.
- InfraTrack only records commands executed through `infratrack run -- ...` while a session is active.
- Captured metadata is minimal: timestamp, sanitized command, exit code, duration, and optional working directory.
- Stdout and stderr are never stored in MVP.
- Redaction happens before writing to disk.
- Denylisted commands are stored as `[REDACTED BY POLICY]`.
- InfraTrack does not perform telemetry, analytics, or network calls in MVP.

## Tests

Run all tests:

```bash
go test ./...
```
