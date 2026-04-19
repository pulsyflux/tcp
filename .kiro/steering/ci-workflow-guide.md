# CI Workflow Guide

## Rules

- The GitHub Actions CI workflow MUST be identical regardless of where it runs — GitHub Actions, local via Act, or any other runner. No conditional logic, no fallbacks, no environment-specific branches.
- The workflow MUST run successfully locally using `act` with `-self-hosted` mode on Windows before being pushed to GitHub.
- Use a single `build-and-test` job for Go vet, test, and vulnerability check.
- The publish job is separate because it only runs on pushes to main (not PRs) and needs to create/push tags.
- Always test workflow changes locally with `run-ci.ps1` before committing.

## Security Order

The build-and-test job MUST run steps in this order:

1. `go mod download` — fetch dependencies
2. `go mod verify` — verify checksums against `go.sum` (fails if tampered)
3. `govulncheck` — check for known vulnerabilities in dependencies
4. `go vet` — static analysis (first step that executes dependency code)
5. `go test` — run tests

Dependencies MUST be verified and vulnerability-checked BEFORE any step that compiles or executes code referencing them (`vet`, `test`, `build`).

The workflow also sets:
- `GOFLAGS: -mod=readonly` — prevents `go.mod`/`go.sum` from being modified during CI
- `GONOSUMDB: ""` — ensures all modules are verified against the Go checksum database

## Failure Behavior

- `go mod verify` — exits non-zero if any downloaded module doesn't match its `go.sum` checksum. Catches supply-chain tampering.
- `govulncheck` — exits with code 3 if vulnerabilities are found that your code actually *calls*. Vulnerabilities in imported packages that your code doesn't call (e.g. unused stdlib paths) do NOT fail the build. This is intentional — it avoids false positives from unreachable code.
- GitHub Actions treats any non-zero exit code as a step failure, which blocks the pipeline and prevents `vet`/`test` from running against unverified or vulnerable dependencies.

## Publishing

On every push to `main` (after tests pass), the `publish` job automatically:

1. Finds the latest `v*` tag (starts at `v0.1.0` if none exists)
2. Bumps the patch version (e.g. `v0.1.0` → `v0.1.1`)
3. Creates and pushes the new git tag
4. Pings the Go module proxy (`proxy.golang.org`) via `go list -m` to cache the version

For minor or major bumps, manually push a tag (e.g. `v0.2.0`, `v1.0.0`) — the workflow will auto-increment from there.

### pkg.go.dev Indexing

- The Go module proxy caches the version immediately — `go get` works right away.
- The `pkg.go.dev` search index is separate and can take up to 24 hours for new modules to appear in search results.
- Direct links work immediately: `https://pkg.go.dev/github.com/pulsyflux/tcp`
- This is normal behavior for first-time modules — no action needed, search catches up on its own.
