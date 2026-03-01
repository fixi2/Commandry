# Security Policy

## Supported Versions

Commandry currently supports security fixes for:

- the latest released version
- the current `main` branch while preparing the next release

Older releases may not receive fixes.

## Reporting a Vulnerability

Use GitHub private vulnerability reporting (Security Advisories) for responsible disclosure.

If private reporting is unavailable, open a public issue **without exploit details** and request a private contact channel.

Do **not** include any of the following in a public issue:

- secrets, tokens, credentials, or private infrastructure details
- proof-of-concept exploit code
- precise exploit steps that would help third parties reproduce the issue

When reporting, include:

- Commandry version (`cmdry version`)
- OS and shell
- install method (`setup`, release binary, local build, etc.)
- clear reproduction steps
- expected vs actual behavior
- impact summary
- relevant output or runbook snippets with secrets removed

## Response Expectations

We aim to:

- acknowledge the report quickly
- confirm whether we can reproduce it
- keep the reporter informed about remediation status

Please use the normal issue templates for bugs, UX feedback, documentation problems, and feature requests. Public issues are not the right place for private security disclosures.