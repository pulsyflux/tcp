# CI Workflow Guide

## Rules

- The GitHub Actions CI workflow MUST be identical regardless of where it runs — GitHub Actions, local via Act, or any other runner. No conditional logic, no fallbacks, no environment-specific branches.
- The workflow MUST run successfully locally using `act` with `-self-hosted` mode on Windows before being pushed to GitHub.
- Use a single `build-and-test` job for Go vet, test, and vulnerability check.
- Always test workflow changes locally with `run-ci.ps1` before committing.
