# Contributing

Thanks for helping improve Commandry.

Commandry is a local-first CLI for recording shell sessions and exporting deterministic Markdown runbooks. The project is still early, so small, focused, well-tested changes are much easier to review and safer to merge than broad mixed changes.

## Before you open an issue or PR

- Read the main usage flow in `README.md`.
- Check `TESTING.md` for the current test packs.
- Use the matching issue template for bugs, UX feedback, documentation issues, feature requests, or beta feedback.
- For non-trivial changes, open or link an issue first so scope is clear before implementation starts.

## Rebrand Notes

The project was renamed from InfraTrack to Commandry during the `v0.6.0` rebrand.

- Treat `InfraTrack` / `infratrack` as legacy compatibility text, not as a second product name.
- Do not introduce new legacy-name references in new changes unless compatibility explicitly requires it.
- When touching a legacy file, move it toward `Commandry` where the change is low-risk, or open a focused issue/PR for the remaining cleanup.
- Be deliberate around risky rename areas: public API/contracts, external integrations, migrations, release metadata, and compatibility paths.

## Pull Request Expectations

Keep PRs small, reviewable, and easy to validate.

A good PR should:

- solve one logical problem
- explain the user or reliability impact clearly
- include tests when behavior changes
- update docs when the user flow changes
- avoid unrelated cleanup unless it is inseparable from the change

If your change affects setup, hooks, export, or CLI output, call that out explicitly in the PR description.

## Commit Message Style

Default format:

`<type>(<scope>): <imperative summary>`

Examples:

- `feat(cli): add cmdry root alias and help updates`
- `fix(setup): migrate legacy config dir to commandry`
- `test(export): update runbook golden for branding changes`
- `chore(ci): pin github actions by commit sha`

### Rules

- Use English.
- Keep the summary short and specific (target: <= 72 chars).
- Do not include version tags in commit titles (for example `v0.6.0`).
- One commit should represent one logical change.

### Recommended types

- `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`

### Recommended scopes

- `cli`, `setup`, `hooks`, `export`, `policy`, `store`, `ci`, `tests`, `docs`, `rebrand`

### Controlled exceptions

This style is the default and should be followed in normal work.

Exceptions are allowed when needed:

- a better type/scope exists for the specific change but is not listed above
- a temporary or emergency change needs a narrower custom scope

When using an exception:

- keep the same overall format
- keep wording clear and concrete
- avoid inventing many new types/scopes without reason

## Merge Rules

- Use fast-forward merge when the branch already has a clean, intentional history.
- Use a non-fast-forward merge only when keeping the branch boundary is useful as a distinct milestone.
- If a branch has noisy or duplicate commits, clean the branch history before merging.

## Branch Hygiene

- Keep one branch focused on one logical change.
- Do not mix code changes, docs cleanup, and unrelated test churn in one branch unless they are inseparable.
- Before commit/merge, remove temporary local artifacts (`*.exe`, `*.tmp`, `*.go.<digits>`, local transcripts) that should not enter the repository.

## Testing Expectations

Before asking for review:

- run the smallest relevant test pack for your change
- include exact commands when a manual flow matters
- mention any environment-specific limitations (for example local antivirus interfering with `go test`)

If you could not run a relevant test, say that explicitly.

## Optional local commit template

This repository includes `.gitmessage` to help keep commit titles consistent.

Use it locally:

`git config commit.template .gitmessage`