# Contributing

## Commit Message Style

Default format:

`<type>(<scope>): <imperative summary>`

Examples:

- `feat(cli): add cmdry root alias and help updates`
- `fix(setup): migrate legacy infratrack config dir to commandry`
- `test(export): update runbook golden for branding changes`
- `chore(ci): pin github actions by commit sha`

### Rules

- Use English.
- Keep summary short and specific (target: <= 72 chars).
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

## Optional local commit template

This repository includes `.gitmessage` to help keep commit titles consistent.

Use it locally:

`git config commit.template .gitmessage`

