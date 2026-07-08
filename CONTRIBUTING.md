# Contributing to Bonus Service

Thank you for contributing to this project. This guide covers the workflow, setup, and review expectations for human contributors.

See [`README.md`](README.md) for a project overview and [`AGENTS.md`](AGENTS.md) for agent-specific guidance.

---

## Setup

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- `golangci-lint` (for local linting)
- `govulncheck` (optional, run in CI)

### Local development

```bash
cd bonus_service
cp .env.example .env
# Edit .env if needed (defaults work for Docker Compose)
docker compose up --build -d
```

The service will be available at `http://localhost:8080/`.

### Run without Docker

```bash
cd bonus_service
export DB_DSN="postgres://admin:admin@localhost:5432/admindb?sslmode=disable"
go run ./cmd/api
```

Apply migrations manually if running without Docker Compose:

```bash
migrate -path ./migrations -database "$DB_DSN" up
```

---

## Verification before submitting

Run these checks locally before opening a Pull Request:

```bash
# Build
go build ./...

# Lint
golangci-lint run ./...

# Unit and integration tests
go test ./internal/... -v
go test ./tests/integration/... -v

# Performance/QRT tests (optional locally)
go test -run=TestPerformance ./... -bench=. -benchtime=10s
```

All of the above run automatically in CI on every PR. Do not open a PR you know will fail CI.

---

## Branch and PR workflow

1. Open or find a GitHub Issue for the work you plan to do.
2. Create a branch from `main` named `<issue-number>-short-description`  
   (e.g. `241-fix-cancel-hold-bug`).
3. Make your changes and commit them to that branch.
4. Open a Pull Request linked to the Issue.
5. Fill in the PR template, including the changelog checklist (see below).
6. Wait for CI to pass and at least one team member to approve.
7. After approval and green CI, merge using **Merge Commit** strategy.

Direct pushes to `main` are disabled. All changes must go through a PR.

---

## Pull Request Template

When opening a Pull Request, fill in the following:

```markdown
# Pull Request
## Related Issue
Closes #[issue number]

## Type of Change
- [ ] Bug fix
- [ ] New feature (non-breaking)
- [ ] Breaking change
- [ ] Database migration
- [ ] Documentation update
- [ ] Infrastructure change

## Description
<!-- Brief description of what this PR does -->

## Acceptance Criteria Verification
- [ ] All acceptance criteria from the issue are satisfied
- [ ] Code is reviewed by another team member
- [ ] CHANGELOG.md is updated (if user-visible change)
- [ ] Not applicable: no user-visible change

## Testing (if applicable)
<!-- Describe how this PR was tested -->

## Screenshots (if applicable)
<!-- Add screenshots here -->
```

For user-visible changes, add an entry under `## [Unreleased]` in `CHANGELOG.md` using the appropriate category: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, or `Security`.

---

## Review expectations

The reviewer checks:

- Acceptance criteria from the linked Issue are satisfied
- Code style and formatting match the existing codebase
- Appropriate tests are included or updated
- No secrets, credentials, or sensitive data are committed
- `CHANGELOG.md` is updated when required
- Documentation is updated when behaviour or setup changes

PR authors cannot approve their own changes. At least one approval is required before merge.

---

## Useful documentation

- [`docs/development-process.md`](docs/development-process.md) — full workflow, CI/CD, and deployment details
- [`docs/definition-of-done.md`](docs/definition-of-done.md) — Definition of Done checklist
- [`docs/testing.md`](docs/testing.md) — testing strategy
- [`docs/migrations.md`](docs/migrations.md) — database schema and migration guide
- [`docs/quality-requirements.md`](docs/quality-requirements.md) — quality gates
