<h1 align="center">Commandry</h1>

<p align="center">
  <img src="docs/assets/commandry-mark.svg" alt="Commandry logo" width="96">
</p>

<p align="center">Local-first CLI for recording shell sessions and exporting deterministic Markdown runbooks.</p>

<p align="center">
  <a href="https://github.com/fixi2/Commandry/releases"><img src="https://img.shields.io/badge/status-beta-ff5a36?style=for-the-badge" alt="Beta"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22+"></a>
  <img src="https://img.shields.io/badge/local--first-0ea5a8?style=for-the-badge" alt="Local-first">
  <img src="https://img.shields.io/badge/telemetry-none-1f8f5f?style=for-the-badge" alt="No telemetry">
</p>

<p align="center">
  <a href="https://github.com/fixi2/Commandry/releases"><img src="https://img.shields.io/badge/download-latest%20release-111827?style=for-the-badge&logo=github" alt="Download latest release"></a>
  <a href="https://github.com/fixi2/Commandry/issues/new?template=beta_feedback.yml"><img src="https://img.shields.io/badge/beta-feedback-111827?style=for-the-badge&logo=github" alt="Beta feedback"></a>
  <a href="TESTING.md"><img src="https://img.shields.io/badge/tests-testing%20guide-111827?style=for-the-badge" alt="Testing guide"></a>
</p>

> Commandry is currently in beta. The project was renamed from InfraTrack in `v0.6.0`; a few migration-safe technical references still keep the legacy name while the remaining contracts finish moving.

## Why Commandry

- Capture real shell work as explicit, reviewable steps.
- Export deterministic Markdown runbooks that are easy to re-read and share.
- Keep data local under your user config directory.
- Stay in control: no telemetry, no background service, no hidden sync.

## Install

Recommended path (pre-`v1.0.0`):

1. Download the latest binary from [GitHub Releases](https://github.com/fixi2/Commandry/releases).
2. Run setup from the downloaded file path once:
   - Windows: `.\cmdry.exe setup`
   - Linux/macOS: `./cmdry setup`
3. Open a new terminal.
4. Verify the install: `cmdry setup status`

What setup does:

- `cmdry setup` - interactive install flow (install binary + update user PATH)
- `cmdry setup plan` - preview setup actions without applying changes
- `cmdry setup apply --yes` - apply setup directly
- `cmdry setup undo` - revert setup changes recorded in setup state

Updating from a freshly downloaded binary:

- `setup` always installs the **currently running** executable.
- To upgrade, run setup from the new downloaded file first.
- Example:
  - Windows: `.\cmdry.exe setup apply --yes`
  - Linux/macOS: `./cmdry setup apply --yes`
- Then open a new terminal and check `cmdry version`.

## Quick Start

### First run

```bash
cmdry i
cmdry s "Deploy to staging" -e staging
cmdry r -- kubectl apply -f deploy.yaml
cmdry r -- kubectl rollout status deploy/api
cmdry stp
cmdry x -l -f md
```

Output:

- `runbooks/<timestamp>-<slug>.md`

### Already installed and initialized

```bash
cmdry s "my runbook session"
cmdry r -- <your command>
cmdry stp
cmdry x -l -f md
```

## Beta feedback

If you are one of the first public testers, use the beta feedback form so early friction lands in the right queue:

- [Beta feedback](https://github.com/fixi2/Commandry/issues/new?template=beta_feedback.yml) - first-impression reports from real use
- [Bug report](https://github.com/fixi2/Commandry/issues/new?template=bug_report.yml) - reproducible problems or regressions
- [UX / CLI feedback](https://github.com/fixi2/Commandry/issues/new?template=ux_feedback.yml) - confusing output, naming, hints, or workflow friction
- [Documentation issue](https://github.com/fixi2/Commandry/issues/new?template=docs_issue.yml) - outdated or unclear docs

## Optional: Shell Hooks (Faster Workflow)

Hooks make capture feel natural in daily shell usage. Once hooks are installed and enabled, commands typed between `start` and `stop` are captured automatically - you do not need to prefix each command with `cmdry run`.

PowerShell:

```bash
cmdry hooks install powershell --yes
cmdry hooks enable
cmdry hooks status
```

Bash:

```bash
cmdry hooks install bash
cmdry hooks enable
cmdry hooks status
```

Zsh:

```bash
cmdry hooks install zsh
cmdry hooks enable
cmdry hooks status
```

Remove hooks at any time:

```bash
cmdry hooks uninstall powershell
cmdry hooks uninstall bash
cmdry hooks uninstall zsh
```

## Command Guide

### Core flow

- `cmdry init` (`i`) - initialize local config and session storage.
- `cmdry start "<title>"` (`s`) - start a recording session. Optional environment label: `--env` / `-e`.
- `cmdry run -- <cmd ...>` (`r`) - execute a command and record a sanitized step.
- `cmdry stop` (`stp`) - finish the active session.
- `cmdry export --last -f md` (`x`) - export the latest completed session to Markdown.

### Helpful commands

- `cmdry status` - show current recording state.
- `cmdry doctor` - run local diagnostics (paths, write access, PATH hints, tool availability).
- `cmdry sessions list -n <count>` - list recent completed sessions.
- `cmdry export --session <id> -f md` - export a specific completed session.
- `cmdry alias --shell <powershell|bash|zsh|cmd>` - print alias snippet for `cmdr` without changing system config.
- `cmdry version` (`v`) - print build version metadata.

### Setup commands

- `cmdry setup`
- `cmdry setup plan`
- `cmdry setup apply`
- `cmdry setup status`
- `cmdry setup undo`

## Troubleshooting

- `cmdry doctor` is the first check when PATH, setup, or local storage looks wrong.
- `cmdry setup status` confirms which binary/path is active.
- `cmdry hooks status` confirms hook install state and whether hooks mode is enabled.

<details>
<summary>Windows shell builtins</summary>

On Windows, commands like `echo`, `dir`, and `copy` are `cmd.exe` builtins, not standalone executables.

Use one of these forms:

```bash
cmdry run -- cmd /c echo "build started"
cmdry run -- powershell -NoProfile -Command "Write-Output 'build started'"
```

</details>

## Data and Privacy

Commandry stores local state under `os.UserConfigDir()/commandry`.

Typical paths:

- Windows: `%APPDATA%\commandry`
- macOS: `~/Library/Application Support/commandry`
- Linux: `~/.config/commandry`

Stored files:

- `config.yaml` - policy and config
- `sessions.jsonl` - completed sessions
- `active_session.json` - current in-progress session (only while recording)

Security defaults:

- Recording is off by default.
- Only explicitly executed commands are captured.
- Captured metadata is minimal: timestamp, sanitized command, exit code, duration, and optional working directory.
- Stdout and stderr are not stored.
- Redaction happens before data is written to disk.
- Denylisted commands are stored as `[REDACTED BY POLICY]` by default.
- Optional: set `policy.enforce_denylist: true` in `config.yaml` to block denylisted commands before execution in `cmdry run`.

Quick examples:

- `cmdry run -- curl -H "Authorization: Bearer abcdef" https://example.com` -> token value is stored as `[REDACTED]`
- `cmdry run -- printenv` -> stored command becomes `[REDACTED BY POLICY]`

Reset / uninstall:

1. Stop any active recording: `cmdry stop`
2. Delete the local `commandry` directory in your user config location
3. Optionally delete project-local `runbooks/`

## Build from Source

Requirements:

- Go 1.22 or newer

```bash
go mod tidy
go build ./cmd/cmdry
```

This creates `cmdry` (or `cmdry.exe` on Windows).

## Tests

Run the full test suite:

```bash
go test ./...
```

Run the black-box contract suite:

```bash
go test ./e2e/blackbox -count=1
```

More test packs and CI notes: [`TESTING.md`](TESTING.md)

## Contributing

- Contribution guide: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- UX behavior contract: [`docs/testing/ux-contract.md`](docs/testing/ux-contract.md)
- Behavior contract: [`docs/testing/behavior-contract.md`](docs/testing/behavior-contract.md)

<details>
<summary>Antivirus note for contributors (Windows)</summary>

Some antivirus products may flag temporary Go test binaries (for example `*.test.exe`) during `go test`. This is usually a false positive for local build artifacts, not Commandry runtime behavior.

If this happens, run tests with local cache/temp folders inside the repo:

```powershell
Set-Location <repo-path>
New-Item -ItemType Directory -Force .gocache, .gotmp | Out-Null
$env:GOCACHE = "$PWD\.gocache"
$env:GOTMPDIR = "$PWD\.gotmp"
go test ./...
```

If your AV still blocks test binaries, add a narrow exclusion for `.gocache` and `.gotmp` instead of disabling protection globally.

</details>

