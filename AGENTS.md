# AGENTS.md

Agent-facing guide for coding agents working in this repository.  
Human contributor guidance is in [`CONTRIBUTING.md`](CONTRIBUTING.md). Project overview is in [`README.md`](README.md).

---

## Repository layout

```bonus_service/     Go service source (cmd/, internal/, migrations/, tests/)
docs/              Project documentation
reports/           Weekly sprint reports
api/               OpenAPI spec and Postman collection
```

All application code lives under `bonus_service/`. Run all Go commands from that directory.

---

## Setup and build

If you build it locally, do next instructions:

```bash
cd bonus_service
cp .env.example .env          # defaults work with Docker Compose
docker compose up --build -d  # starts API + PostgreSQL
```

If you build it on VM, you should change server in api/openapi.yaml:

```bash
servers:
  - url: http://your_server_ip:8080
    description: Some discription
```

and then do previous step

Without Docker:

```bash
cd bonus_service
export DB_DSN="postgres://admin:admin@localhost:5432/admindb?sslmode=disable"
migrate -path ./migrations -database "$DB_DSN" up
go run ./cmd/api
```

---

## Verification commands

Run all of these before committing or opening a PR:

```bash
cd bonus_service

# Build
go build ./...

# Lint
golangci-lint run ./...

# Unit and integration tests
go test ./internal/... -v
go test ./tests/integration/... -v

# Vulnerability scan
govulncheck ./...
```

CI runs the same checks. Do not open PRs that fail these commands.

---

## Workflow

- `main` is the protected default branch. Do **not** push directly to it.
- Branch naming: `<issue-number>-short-description` (e.g. `241-fix-cancel-hold`).
- All changes go through a Pull Request linked to a GitHub Issue.
- At least one human approval is required before merging.
- Use **Merge Commit** as the merge strategy.
- Every PR must include a changelog checklist item (see [`CONTRIBUTING.md`](CONTRIBUTING.md)).

---

## Changelog

When a change is user-visible, add an entry under `## [Unreleased]` in `CHANGELOG.md` using one of: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`.

If the change is not user-visible, note that in the PR description instead.

---

## Safety and credential handling

- **Never commit secrets, credentials, or `.env` files.** `.env` is in `.gitignore`.
- Environment variables are provided via `.env` (local) or CI secrets — never hardcoded.
- Required variables: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `DB_DSN`.
- Do not modify or expose the deployed instance at `http://10.93.26.189:8080/`.
- Do not run destructive migrations (`down`) against shared or production databases.
- Do not move, replace, or delete Git tags mapped to submitted MVP milestones (`v1.0.0`, `v1.1.0`, `v2.0.0`).

---

## Deeper documentation

- [`docs/development-process.md`](docs/development-process.md) — full workflow and CI/CD details
- [`docs/architecture/README.md`](docs/architecture/README.md) — architecture overview and ADRs
- [`docs/migrations.md`](docs/migrations.md) — database schema and migration guide
- [`docs/testing.md`](docs/testing.md) — testing strategy
- [`docs/definition-of-done.md`](docs/definition-of-done.md) — Definition of Done
- [`CHANGELOG.md`](CHANGELOG.md) — release history
